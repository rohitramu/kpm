package common

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/emirpasic/gods/stacks/linkedliststack"

	"github.com/rohitramu/kpm/subcommands/utils/constants"
	"github.com/rohitramu/kpm/subcommands/utils/log"
	"github.com/rohitramu/kpm/subcommands/utils/templates"
	"github.com/rohitramu/kpm/subcommands/utils/types"
	"github.com/rohitramu/kpm/subcommands/utils/validation"
	"github.com/rohitramu/kpm/subcommands/utils/yaml"
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
	hash              *string

	OutputName          string
	PackageDirPath      string
	ExecutableTemplates []*template.Template
	TemplateInput       *types.GenericMap
}

// VisitNodesDepthFirst visits nodes in the tree in depth-first fashion, applying the given consumer function on each node.  It returns the number of nodes that were visited.
func (tree *DependencyTree) VisitNodesDepthFirst(consumeNode func(relativeFilePath []string, friendlyNamePath []string, executableTemplates []*template.Template, templateInput *types.GenericMap) error) (int, error) {
	var err error
	var ok bool

	var numVisitedNodes = 0
	var toVisitStack = linkedliststack.New()
	toVisitStack.Push(tree.root)
	for !toVisitStack.Empty() {
		// Get the next item to visit off of the stack
		var nodeObj interface{}
		if nodeObj, ok = toVisitStack.Pop(); !ok {
			log.Panic("Failed to get next node")
		}

		// Convert the object into a node
		var node *dependencyTreeNode
		node, ok = nodeObj.(*dependencyTreeNode)
		if !ok {
			log.Panic("Failed to cast item in stack to node object")
		}

		// Keep track of the number of visited nodes
		numVisitedNodes++

		// Get the current path excluding the root (use stacks so it is easy to reverse the list later)
		var friendlyNameStack = linkedliststack.New()
		var relativeFileStack = linkedliststack.New()
		var currentNode = node
		for currentNode != nil {
			// Get the package info
			var packageInfo = currentNode.packageDefinition.PackageInfo
			var packageFullName = constants.GetPackageFullName(packageInfo.Name, packageInfo.Version)

			// Get the next segment in the friendly name path
			var friendlyName = GetOutputFriendlyName(currentNode.OutputName, packageFullName)

			// Append the next friendly name path segment to the friendly name path
			friendlyNameStack.Push(friendlyName)

			// Append the output name to the relative file path
			relativeFileStack.Push(currentNode.OutputName)

			// Update the current node
			currentNode = currentNode.Parent
		}

		// Get the reversed list so it is ordered from the top (the root) to the bottom of the tree
		var friendlyNamePath = make([]string, friendlyNameStack.Size())
		var relativeFilePath = make([]string, relativeFileStack.Size())
		var friendlyNameIterator = friendlyNameStack.Iterator()
		var relativeFileIterator = relativeFileStack.Iterator()
		var i = 0
		for friendlyNameIterator.Next() {
			var val = friendlyNameIterator.Value()
			friendlyNamePath[i], ok = val.(string)
			if !ok {
				log.Panic("Unexpected object type while iterating over friendly name path segments: %s", reflect.TypeOf(val))
			}
			i++
		}
		i = 0
		for relativeFileIterator.Next() {
			var val = relativeFileIterator.Value()
			relativeFilePath[i], ok = val.(string)
			if !ok {
				log.Panic("Unexpected object type while iterating over relative file path segments: %s", reflect.TypeOf(val))
			}
			i++
		}

		// Call the consuming function
		if err = consumeNode(relativeFilePath, friendlyNamePath, node.ExecutableTemplates, node.TemplateInput); err != nil {
			return 0, err
		}

		// Visit all the children
		for _, childNode := range node.Children {
			toVisitStack.Push(childNode)
		}
	}

	return numVisitedNodes, nil
}

// GetDependencyTree ensures that the dependency tree has no loops and then returns the dependency tree.
func GetDependencyTree(kpmHomeDir string, packageName string, packageVersion string, dockerRegistry string, outputName string, parameters *types.GenericMap) (*DependencyTree, error) {
	var err error
	var ok bool

	// Validate output name
	if err = validation.ValidateOutputName(outputName); err != nil {
		return nil, fmt.Errorf("Invalid output name \"%s\" for package: %s\n%s", outputName, packageName, err)
	}

	// Create the package definition for the root node
	var rootNodePackageDefinition = &types.PackageDefinition{
		PackageInfo: &types.PackageInfo{
			Name:    packageName,
			Version: packageVersion,
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
		var currentNode *dependencyTreeNode
		currentNode, ok = currentNodeObj.(*dependencyTreeNode)
		if !ok {
			log.Panic("Object on \"toVisit\" list is not a tree node")
		}

		log.Verbose("Visiting node: %s", currentNode.OutputName)

		var currentOutputName = currentNode.OutputName
		var currentPackageName = currentNode.packageDefinition.PackageInfo.Name
		var currentPackageVersion = currentNode.packageDefinition.PackageInfo.Version
		var currentParameters = currentNode.packageDefinition.Parameters

		// Validate package name
		err = validation.ValidatePackageName(currentPackageName)
		if err != nil {
			return nil, fmt.Errorf("Invalid name for package \"%s\": %s", currentOutputName, err)
		}

		// Validate package version
		err = validation.ValidatePackageVersion(currentPackageVersion)
		if err != nil {
			return nil, fmt.Errorf("Invalid version for package \"%s\": %s", currentOutputName, err)
		}

		// Get the package's full name
		var currentPackageFullName = constants.GetPackageFullName(currentPackageName, currentPackageVersion)

		// Make sure that the parameters were provided
		if currentParameters == nil {
			var friendlyName = GetOutputFriendlyName(currentOutputName, currentPackageFullName)
			return nil, fmt.Errorf("Output was not provided any parameters: %s", friendlyName)
		}

		// Check local repository for package
		var packages []string
		packages, err = GetPackageFullNamesFromLocalRepository(kpmHomeDir)
		if err != nil {
			return nil, err
		}
		var found = false
		for _, value := range packages {
			if value == currentPackageFullName {
				found = true
				break
			}
		}
		if !found {
			log.Warning("Package \"%s\" not found in local repository, now checking docker registry \"%s\"...", currentPackageFullName, dockerRegistry)

			// Check remote repository for package
			err = PullPackage(kpmHomeDir, dockerRegistry, currentPackageName, currentPackageVersion)
			if err != nil {
				return nil, fmt.Errorf("Failed to get package \"%s\" from docker registry \"%s\": %s", currentPackageFullName, dockerRegistry, err)
			}
		}

		// Create a function to easily get the human readable path
		var getFriendlyPath = func() string {
			var segments = make([]string, currentPathNodes.Size()+1)
			var it = currentPathNodes.Iterator()
			var i = 0
			for it.Next() {
				var val = it.Value()
				var nodeVal *dependencyTreeNode
				nodeVal, ok = val.(*dependencyTreeNode)
				if !ok {
					log.Panic("Unexpected type in nodes path: %s", reflect.TypeOf(val))
				}

				var packageDefVal = nodeVal.packageDefinition
				if packageDefVal == nil {
					log.Panic("Unexpected nil value for package definition: %s", nodeVal.PackageDirPath)
				}

				var packageInfoVal = packageDefVal.PackageInfo
				if packageInfoVal == nil {
					log.Panic("Unexpected nil value for package info: %s", nodeVal.PackageDirPath)
				}

				segments[i] = GetOutputFriendlyName(nodeVal.OutputName, constants.GetPackageFullName(packageInfoVal.Name, packageInfoVal.Version))

				i++
			}

			// Add current node
			segments[len(segments)-1] = GetOutputFriendlyName(currentNode.OutputName, currentPackageFullName)

			return strings.Join(segments, " -> ")
		}

		// Get the package directory
		var currentPackageDirPath = constants.GetPackageDir(kpmHomeDir, currentPackageFullName)

		// Create shared template (with common options, functions and helper templates for this package)
		var sharedTemplate *template.Template
		sharedTemplate, err = GetSharedTemplate(currentPackageDirPath)
		if err != nil {
			return nil, fmt.Errorf("Failed to construct shared template in package: %s\n%s", getFriendlyPath(), err)
		}

		// Calculate values to be used as inputs to the templates in this package
		var templateInput *types.GenericMap
		templateInput, err = GetTemplateInput(kpmHomeDir, currentPackageFullName, sharedTemplate, currentParameters)
		if err != nil {
			return nil, fmt.Errorf("Failed to get template input in package: %s\n%s", getFriendlyPath(), err)
		}

		// Get the dependency definition templates
		var dependencyTemplates []*template.Template
		dependencyTemplates, err = GetDependencyDefinitionTemplates(sharedTemplate, currentPackageDirPath)
		if err != nil {
			return nil, fmt.Errorf("Failed to get dependency definition templates in package: %s\n%s", getFriendlyPath(), err)
		}

		// Save the package directory path, shared template and calculated values that can be used with this package in the node
		currentNode.PackageDirPath = currentPackageDirPath
		currentNode.TemplateInput = templateInput
		currentNode.ExecutableTemplates, err = GetExecutableTemplates(sharedTemplate, currentPackageDirPath)
		if err != nil {
			return nil, fmt.Errorf("Failed to get executable templates in package: %s\n%s", getFriendlyPath(), err)
		}

		// Check if there is a loop in the dependency tree
		var currentNodeHash = currentNode.getPackageNodeHash()
		if _, exists := currentPathNodes.Get(currentNodeHash); exists {
			// Found a loop
			var dependencyLoop = make([]string, currentPathNodes.Size()+1)
			for i, keyObj := range currentPathNodes.Keys() {
				if valueObj, found := currentPathNodes.Get(keyObj); !found {
					log.Panic("Failed to find value in path nodes map for key: %s", keyObj)
				} else {
					// Value is the node object
					if value, ok := valueObj.(*dependencyTreeNode); !ok {
						log.Panic("Found value in path nodes map which is not a node, for key: %s", keyObj)
					} else {
						var dependencyPackageName = value.packageDefinition.PackageInfo.Name
						var dependencyPackageVersion = value.packageDefinition.PackageInfo.Version
						var dependencyPackageFullName = constants.GetPackageFullName(dependencyPackageName, dependencyPackageVersion)
						dependencyLoop[i] = GetOutputFriendlyName(value.OutputName, dependencyPackageFullName)

						// Add a special symbol to identify the package causing the problem
						if value.getPackageNodeHash() == currentNodeHash {
							dependencyLoop[i] += " [START]"
						}
					}
				}
			}

			// Add the current node
			dependencyLoop[len(dependencyLoop)-1] = GetOutputFriendlyName(currentOutputName, currentPackageFullName) + " [END]"

			// Return an error with the formatted package path
			return nil, fmt.Errorf("Found a circular reference in the dependency tree:\n%s", strings.Join(dependencyLoop, " -> "))
		}

		// Add this node to the map which is tracking the current path
		currentPathNodes.Put(currentNodeHash, currentNode)

		// Evaluate dependencies
		if len(dependencyTemplates) == 0 {
			// If this node has no children, remove it from the current path
			currentPathNodes.Remove(currentPackageFullName)
		} else {
			// Execute the dependency definition templates to get the concrete dependency definitions
			for _, dependencyTemplate := range dependencyTemplates {
				// Get the dependency template's file name
				var templateFileName = dependencyTemplate.Name()

				// Remove the file extension to get the dependency's output name
				var dependencyOutputName = strings.TrimSuffix(templateFileName, filepath.Ext(templateFileName))

				// Get the package definition by running the template input through the package definition file
				var dependencyDefinitionBytes []byte
				dependencyDefinitionBytes, err = templates.ExecuteTemplate(dependencyTemplate, currentNode.TemplateInput)
				if err != nil {
					return nil, fmt.Errorf("Failed to execute dependency definition template \"%s\" in package: %s\n%s", templateFileName, getFriendlyPath(), err)
				}

				// Create an object from the package definition
				var dependencyDefinition = new(types.PackageDefinition)
				err = yaml.BytesToObject(dependencyDefinitionBytes, dependencyDefinition)
				if err != nil {
					return nil, err
				}

				// Make sure that the package info object is not nil
				if dependencyDefinition.PackageInfo == nil {
					return nil, fmt.Errorf("Package info was not found for dependency of package \"%s\": %s", currentPackageFullName, dependencyOutputName)
				}

				// Push new dependency node
				var dependencyNode *dependencyTreeNode
				if dependencyNode, err = getPackageNode(currentNode, dependencyDefinition, dependencyOutputName, ""); err != nil {
					return nil, err
				}
				toVisitStack.Push(dependencyNode)
			}
		}

		// Make sure to clean up all of the nodes that will no longer be in the path on the next iteration
		var pathIt = currentPathNodes.Iterator()
		found = false
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
							break
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
	var err error

	// Validate inputs
	if err = validation.ValidateOutputName(outputName); err != nil {
		return nil, fmt.Errorf("Invalid output name: %s\n%s", outputName, err)
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

func (node *dependencyTreeNode) getPackageNodeHash() string {
	var err error

	if node == nil {
		log.Panic("Package node cannot be nil")
	}

	if node.hash != nil {
		return *node.hash
	}

	if node.packageDefinition == nil {
		log.Panic("Package definition cannot be nil")
	}
	if node.TemplateInput == nil {
		log.Panic("Template input cannot be nil")
	}
	if node.packageDefinition.PackageInfo == nil {
		log.Panic("Package info inside package definition cannot be nil")
	}

	var hashedValues = struct {
		PackageInfo   *types.PackageInfo
		PackageInputs *types.GenericMap
	}{
		node.packageDefinition.PackageInfo,
		node.TemplateInput,
	}

	var hash []byte
	hash, err = yaml.ObjectToBytes(&hashedValues)
	if err != nil {
		log.Panic("Invalid object for node: %s", node.OutputName)
	}

	var stringHash = string(hash)
	node.hash = &stringHash

	return stringHash
}
