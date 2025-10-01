/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

// kCmd represents the k command
// var kCmd = &cobra.Command{
// 	Use:   "k",
// 	Short: "visualizing refs and relavant oids",
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		var refPath string
// 		if len(args) == 1 {
// 			refPath = filepath.Join(data.RefBranchPath, args[0])
// 		} else {
// 			refPath = data.RefHEADPath
// 		}
// 		ref, err := data.NewRef(refPath)
// 		if err != nil {
// 			return err
// 		}
// 		if ref.IsSymbolic {
// 			ref, err = ref.ResolveSymbolic(ref.Next)
// 			if err != nil {
// 				return fmt.Errorf("internal error: %w", err)
// 			}
// 		}
// 		oid := ref.Oid
// 		refName := filepath.Base(refPath)
// 		parent, err := commitParent(oid)
// 		if err != nil {
// 			return err
// 		}
// 		if parent == "" {
// 			fmt.Printf("%s -> %s", refName, oid)
// 			return nil
// 		}

// 		rootNode := NewRefNode(refName, oid, nil)
// 		tmp := rootNode
// 		m, err := OidRefMap()
// 		if err != nil {
// 			return err
// 		}
// 		for {
// 			child, err := tmp.AddChild(m)
// 			if err != nil {
// 				return err
// 			}
// 			if child == nil {
// 				break
// 			}
// 			tmp = child
// 		}
// 		fmt.Println(rootNode.PrintTree())
// 		return nil
// 	},
// }

// type RefNode struct {
// 	ref   string
// 	oid   string
// 	child *RefNode
// }

// func NewRefNode(ref string, oid string, child *RefNode) *RefNode {
// 	return &RefNode{
// 		ref:   ref,
// 		oid:   oid,
// 		child: child,
// 	}
// }

// func (n *RefNode) AddChild(oidRefMap map[string]string) (child *RefNode, err error) {
// 	if n.child != nil {
// 		return n.child, nil
// 	}
// 	oid, err := commitParent(n.oid)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if oid == "" {
// 		return nil, nil
// 	}
// 	if ref, ok := oidRefMap[oid]; ok {
// 		n.child = NewRefNode(ref, oid, nil)
// 	}
// 	return n.child, nil
// }

// func (n *RefNode) PrintTree() string {
// 	var buf strings.Builder
// 	fmt.Fprintf(&buf, "%s -> %s\n", n.ref, n.oid)

// 	tmp := n
// 	for {
// 		if tmp.child == nil {
// 			break
// 		}
// 		fmt.Fprintf(&buf, "%s -> %s\n", tmp.child.ref, tmp.child.oid)
// 		tmp = tmp.child
// 	}
// 	if buf.Len() == 0 {
// 		return "no tree"
// 	}
// 	return buf.String()
// }

// // creating {key: oid, value: ref} map
// func OidRefMap() (map[string]string, error) {
// 	m := make(map[string]string)
// 	rootdir := data.RefBranchPath
// 	err := filepath.WalkDir(rootdir, func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if rootdir == path {
// 			return nil
// 		}
// 		if d.IsDir() {
// 			return filepath.SkipDir
// 		}
// 		f, err := os.Open(path)
// 		if err != nil {
// 			return err
// 		}
// 		defer f.Close()
// 		oid, err := io.ReadAll(f)
// 		if err != nil {
// 			return err
// 		}
// 		if _, exists := m[string(oid)]; !exists {
// 			m[string(oid)] = d.Name()
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("OidRefMap: %w", err)
// 	}
// 	return m, nil
// }

// func init() {
// 	rootCmd.AddCommand(kCmd)
// }
