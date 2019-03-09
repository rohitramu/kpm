package common

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/emirpasic/gods/stacks/linkedliststack"

	"../utils/constants"
	"../utils/log"
	"../utils/templates"
	"../utils/types"
	"../utils/validation"
	"../utils/yaml"
)

// DependencyTree is the definition of the package dependency tree.
type DependencyTree struct {
	root *dependencyTreeNode
}

// dependencyTreeNode is the definition of a package in a dependency tree.
type dependencyTreeNode struct {
	Parent   *dependencyTreeNode
	Children []*dependencyTreeNode

	packageDefinition *types.PackageDefinition

	OutputName          string
	PackageDirPath      string
	ExecutableTemplates []*template.Template
	TemplateInput       *types.GenericMap
}

// VisitNodesDepthFirst visits nodes in the tree in depth-first fashion, applying the given consumer function on each node.  It returns the number of nodes that were visited.
func (tree *DependencyTree) VisitNodesDepthFirst(consumeNode func(path []string, executableTemplates []*template.Template, templateInput *types.GenericMap) error) (int, error) {
	var numVisitedNodes = 0
	var toVisitStack = linkedliststack.New()
	toVisitStack.Push(tree.root)
	for !toVisitStack.Empty() {
		if nodeObj, ok := toVisitStack.Pop(); !ok {
			log.Panic("Failed to get next node")
		} else {
			if node, ok := nodeObj.(*dependencyTreeNode); !ok {
				log.Panic("Failed to cast item in stack to node object")
			} else {
				numVisitedNodes++

				// Get the current path excluding the root
				var nodePath []string
				var currentNode = node
				for currentNode != nil && currentNode.Parent != nil {
					nodePath = append(nodePath, currentNode.OutputName)
					currentNode = currentNode.Parent
				}

				// Reverse the path so it is ordered from the top (the root) to the bottom of the tree
				for left, right := 0, len(nodePath)-1; left < right; left, right = left+1, right-1 {
					nodePath[left], nodePath[right] = nodePath[right], nodePath[left]
				}

				// Call the consuming function
				if err := consumeNode(nodePath, node.ExecutableTemplates, node.TemplateInput); err != nil {
					return 0, err
				}

				for _, childNode := range node.Children {
					toVisitStack.Push(childNode)
				}
			}
		}
	}

	return numVisitedNodes, nil
}

// GetDependencyTree ensures that the dependency tree has no loops and then returns the dependency tree.
func GetDependencyTree(outputName string, kpmHomeDir string, packageName string, wildcardPackageVersion string, parameters *types.GenericMap) (*DependencyTree, error) {
	var err error

	// Validate output name
	if err := validation.ValidateOutputName(outputName); err != nil {
		return nil, err
	}

	// Get the package repository location
	var packageRepositoryDirPath = constants.GetPackageRepositoryDirPath(kpmHomeDir)

	// Create the package definition for the root node
	var rootNodePackageDefinition = &types.PackageDefinition{
		Package: &types.PackageInfo{
			Name:    packageName,
			Version: wildcardPackageVersion,
		},
		Parameters: parameters,
	}

	// Get the root node
	var rootNode *dependencyTreeNode
	if rootNode, err = getPackageNode(nil, rootNodePackageDefinition, outputName, ""); err != nil {
		return nil, err
	}

	// Create the tree
	var tree = &DependencyTree{
		root: rootNode,
	}

	// Traverse and build the tree in a depth-first fashion, looking for loops
	var currentPathNodes = linkedhashmap.New()
	var toVisitStack = linkedliststack.New()
	toVisitStack.Push(rootNode)
	var i = 0
	for currentNodeObj, notEmpty := toVisitStack.Pop(); notEmpty; currentNodeObj, notEmpty = toVisitStack.Pop() {
		if currentNode, ok := currentNodeObj.(*dependencyTreeNode); !ok {
			log.Panic("Object on \"toVisit\" list is not a tree node")
		} else {
			log.Verbose("Visiting node: %s", currentNode.OutputName)

			outputName = currentNode.OutputName
			packageName = currentNode.packageDefinition.Package.Name
			wildcardPackageVersion = currentNode.packageDefinition.Package.Version
			parameters = currentNode.packageDefinition.Parameters

			// Validate package name
			err = validation.ValidatePackageName(packageName)
			if err != nil {
				return nil, fmt.Errorf("Invalid name for package \"%s\": %s", currentNode.OutputName, err)
			}

			// Validate package version
			err = validation.ValidatePackageVersion(wildcardPackageVersion, true)
			if err != nil {
				return nil, fmt.Errorf("Invalid version for package \"%s\": %s", currentNode.OutputName, err)
			}

			// // Check remote repository for newest matching versions of the package
			// if pulledVersion, err := PullPackage(packageName, wildcardPackageVersion); err != nil {
			// 	log.Warning(err)
			// } else {
			// 	wildcardPackageVersion = pulledVersion
			// }

			// Resolve the package version
			var resolvedPackageVersion string
			if resolvedPackageVersion, err = ResolvePackageVersion(kpmHomeDir, packageName, wildcardPackageVersion); err != nil {
				return nil, err
			}
			// Get the package's full name
			var packageFullName = constants.GetPackageFullName(packageName, resolvedPackageVersion)

			// Check if there is a loop in the dependency tree
			if _, exists := currentPathNodes.Get(packageFullName); exists {
				var dependencyLoop = make([]string, currentPathNodes.Size()+1)
				for i, keyObj := range currentPathNodes.Keys() {
					if valueObj, found := currentPathNodes.Get(keyObj); !found {
						log.Panic("Failed to find value in path nodes map for key: %s", keyObj)
					} else {
						if value, ok := valueObj.(*dependencyTreeNode); !ok {
							log.Panic("Found value in path nodes map which is not a node")
						} else {
							if key, ok := keyObj.(string); !ok {
								log.Panic("Found key in path nodes map which is not a string")
							} else {
								dependencyLoop[i] = fmt.Sprintf("%s (%s)", key, value.OutputName)
							}
						}
					}
				}
				dependencyLoop[len(dependencyLoop)-1] = fmt.Sprintf("%s (%s)", packageFullName, outputName)
				return nil, fmt.Errorf("Found loop in dependency tree: %s", strings.Join(dependencyLoop, "->"))
			}

			// Add this node to the map which is tracking the current path
			currentPathNodes.Put(packageFullName, currentNode)

			// Get the package directory
			var packageDirPath = constants.GetPackageDirPath(packageRepositoryDirPath, packageFullName)

			// Create shared template (with common options, functions and helper templates for this package)
			var sharedTemplate *template.Template
			sharedTemplate, err = GetSharedTemplate(packageDirPath)
			if err != nil {
				return nil, err
			}

			// Calculate values to be used as inputs to the templates in this package
			var templateInput *types.GenericMap
			templateInput, err = GetTemplateInput(sharedTemplate, packageDirPath, parameters)
			if err != nil {
				return nil, err
			}

			// Get the dependency definition templates
			var dependencyTemplates []*template.Template
			dependencyTemplates, err = GetDependencyDefinitionTemplates(sharedTemplate, packageDirPath)
			if err != nil {
				return nil, err
			}

			// Save the package directory path, shared template and calculated values that can be used with this package in the node
			currentNode.PackageDirPath = packageDirPath
			currentNode.TemplateInput = templateInput
			currentNode.ExecutableTemplates, err = GetExecutableTemplates(sharedTemplate, packageDirPath)
			if err != nil {
				return nil, err
			}

			// Evaluate dependencies
			if len(dependencyTemplates) == 0 {
				// If this node has no children, remove it from the current path
				currentPathNodes.Remove(packageFullName)
			} else {
				// Execute the dependency definition templates to get the concrete dependency definitions
				for _, dependencyTemplate := range dependencyTemplates {
					// Get the dependency template's file name
					var templateFileName = dependencyTemplate.Name()

					// Remove the file extension to get the dependency's output name
					var dependencyOutputName = strings.TrimSuffix(templateFileName, filepath.Ext(templateFileName))

					// Get the package definition by running the template input through the package definition file
					var dependencyPackageDefinitionBytes []byte
					dependencyPackageDefinitionBytes, err = templates.ExecuteTemplate(dependencyTemplate, currentNode.TemplateInput)
					if err != nil {
						return nil, err
					}

					// Create an object from the package definition
					var dependencyPackageDefinition = new(types.PackageDefinition)
					err = yaml.BytesToObject(dependencyPackageDefinitionBytes, dependencyPackageDefinition)
					if err != nil {
						return nil, err
					}

					// Push new dependency node
					var dependencyNode *dependencyTreeNode
					if dependencyNode, err = getPackageNode(currentNode, dependencyPackageDefinition, dependencyOutputName, ""); err != nil {
						return nil, err
					}
					toVisitStack.Push(dependencyNode)
				}
			}

			// Make sure to clean up all of the nodes that will no longer be in the path on the next iteration
			var pathIt = currentPathNodes.Iterator()
			var found = false
			for pathIt.End(); pathIt.Prev() && !found; {
				// Get this node's children
				if pathNode, pathNodeIsCorrectType := pathIt.Value().(*dependencyTreeNode); !pathNodeIsCorrectType {
					log.Panic("Path object is not a tree node")
				} else {
					for _, childNode := range pathNode.Children {
						// Check if the "toVisit" stack contains this child node
						var stackIt = toVisitStack.Iterator()
						for stackIt.Begin(); stackIt.Next() && !found; {
							if stackNode, ok := stackIt.Value().(*dependencyTreeNode); !ok {
								log.Panic("Stack object is not a tree node")
							} else if stackNode == childNode {
								found = true
							}
						}
					}
				}
			}
		}

		i++
	}

	return tree, nil
}

func getPackageNode(parentNode *dependencyTreeNode, packageDefinition *types.PackageDefinition, outputName string, packageDirPath string) (*dependencyTreeNode, error) {
	// Validate inputs
	if err := validation.ValidateOutputName(outputName); err != nil {
		return nil, err
	}

	// Create the node
	var packageNode = &dependencyTreeNode{
		Parent:   parentNode,
		Children: []*dependencyTreeNode{},

		packageDefinition: packageDefinition,

		OutputName:     outputName,
		PackageDirPath: packageDirPath,
	}

	// If this is not the root node, add this as a child to the parent node
	if parentNode != nil {
		parentNode.Children = append(parentNode.Children, packageNode)
	}

	return packageNode, nil
}
