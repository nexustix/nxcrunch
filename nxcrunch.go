package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"

	bp "github.com/nexustix/boilerplate"
	nrc "github.com/nxReplicator/nxReplicatorCommon"
)

//nxcrunch generate amazingMolecule amazingBulk

func main() {
	fmt.Printf("Hello World Mk II\n")

	args := os.Args

	usr, err := user.Current()
	bp.FailError(err)
	workingDir := usr.HomeDir
	atomDir := nrc.InitWorkFolder(workingDir, ".nxreplicator", "atoms")
	moleculeDir := nrc.InitWorkFolder(workingDir, ".nxreplicator", "molecules")
	bulkDir := nrc.InitWorkFolder(workingDir, ".nxreplicator", "bulks")

	atomManager := nrc.AtomManager{WorkingDir: atomDir}
	//molecule := nrc.Molecule{}

	action := bp.StringAtIndex(1, args)
	moleculeID := bp.StringAtIndex(2, args)
	bulkID := bp.StringAtIndex(3, args)

	fmt.Printf("%s | %s | %s\n", action, moleculeID, bulkID)

	switch action {
	case "generate":
		if (moleculeID != "") && (bulkID != "") {
			generateBulk(&atomManager, moleculeID, moleculeDir, bulkID, bulkDir)
		}
	}
}

func generateBulk(atomManager *nrc.AtomManager, moleculeID, moleculeDir, bulkID, bulkDir string) {
	moleculePath := path.Join(moleculeDir, moleculeID+".nxrm")
	bulkPath := path.Join(bulkDir, bulkID+".nxrb")

	tmpMolecule := nrc.Molecule{}
	tmpMolecule.LoadFromFile(moleculePath)

	tmpBulk := nrc.Bulk{}
	//tmpBulk.LoadFromFile(bulkPath)

	//fmt.Println("-----")
	for k, v := range tmpMolecule.MoleculeItems {
		fmt.Printf("#####(%v of %v)##### %s\n", k+1, len(tmpMolecule.MoleculeItems), v.AtomID)
		if atomManager.HasEntry(v.ProviderID, v.AtomID) {
			//fmt.Printf("<-> fetching '%s'", v.AtomID)

			tmpBulkItem := nrc.BulkItem{}
			tmpDownload := bp.Download{}

			fmt.Printf("</> EXEC >%s %s %s %s<\n", "nxatomize", "downinfo", v.ProviderID, v.AtomID)

			providerCommand := exec.Command("nxatomize", "downinfo", v.ProviderID, v.AtomID)
			output, err := providerCommand.Output()
			if bp.GotError(err) {
				fmt.Printf("<!> ERROR getting downinfo of: '%s'\n", v.AtomID)
			} else {
				//fmt.Printf("{%s}", output)
				fmt.Printf("<-> done getting downinfo of: '%s'\n", v.AtomID)
			}
			lines := strings.Split(string(output), "\n")
			lastLine := bp.StringAtIndex(len(lines)-2, lines) // -2 since output ends with the seperator
			//lastLine := lines[0]
			if (lastLine != "") && !strings.HasPrefix(lastLine, "<") {
				fmt.Printf("<~> %s\n", lastLine)
				downinfos := strings.SplitN(lastLine, "|", 2)
				if len(downinfos) >= 2 {
					tmpDownload.URL = downinfos[0]
					tmpDownload.Filename = downinfos[1]

					tmpBulkItem.RelativePath = v.Dir
					tmpBulkItem.Download = tmpDownload
					//tmpBulk.AddDownload(tmpDownload)
					tmpBulk.AddDownload(tmpBulkItem)
				} else {
					fmt.Printf("<!> ERROR downinfo of: '%s' corrupt\n", v.AtomID)
				}
			}
		} else {
			fmt.Printf("<!> WARNING unable to resolve Atom '%s' via '%s' Provider\n", v.AtomID, v.ProviderID)
			//TODO do wildcard search ?
		}
		//fmt.Println("-----")
	}
	tmpBulk.SaveToFile(bulkPath)

}
