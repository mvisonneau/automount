package lvm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mvisonneau/automount/pkg/exec"
)

// PhysicalVolume represents a LVM physical volume
type PhysicalVolume struct {
	Name        string
	SizeInBytes int
	FreeInBytes int
	*VolumeGroup
}

// VolumeGroup represents a LVM volume group
type VolumeGroup struct {
	Name        string
	SizeInBytes int
	FreeInBytes int
	Tags        []string
	LogicalVolumes
	PhysicalVolumes
}

// LogicalVolume represents a LVM logical volume
type LogicalVolume struct {
	Name        string
	SizeInBytes int
	Tags        []string
	*VolumeGroup
}

// PhysicalVolumes represents a slice of physical volumes
type PhysicalVolumes []*PhysicalVolume

// VolumeGroups represents a slice of volume groups
type VolumeGroups []*VolumeGroup

// LogicalVolumes represents a slice of logical volumes
type LogicalVolumes []*LogicalVolume

// LVM is used for supporting maps of PVs, VGs and LVs
type LVM struct {
	pvs map[string]*PhysicalVolume
	vgs map[string]*VolumeGroup
	lvs map[string]*LogicalVolume
}

// PhysicalVolume returns a physical volume object based on its name
func (l *LVM) PhysicalVolume(name string) (*PhysicalVolume, error) {
	if _, exists := l.pvs[name]; !exists {
		return nil, fmt.Errorf("Cannot find physical volume : '%s'", name)
	}

	return l.pvs[name], nil
}

// VolumeGroup returns a volume group object based on its name
func (l *LVM) VolumeGroup(name string) (*VolumeGroup, error) {
	if _, exists := l.vgs[name]; !exists {
		return nil, fmt.Errorf("Cannot find volume group : '%s'", name)
	}

	return l.vgs[name], nil
}

// LogicalVolume returns a logical volume object based on its name
func (l *LVM) LogicalVolume(name string) (*LogicalVolume, error) {
	if _, exists := l.lvs[name]; !exists {
		return nil, fmt.Errorf("Cannot find logical volume : '%s'", name)
	}

	return l.lvs[name], nil
}

// PhysicalVolumes returns a list of PhysicalVolumes
func (l *LVM) PhysicalVolumes() (pvs PhysicalVolumes) {
	for _, pv := range l.pvs {
		pvs = append(pvs, pv)
	}
	return
}

// VolumeGroups returns a list of VolumeGroups
func (l *LVM) VolumeGroups() (vgs VolumeGroups) {
	for _, vg := range l.vgs {
		vgs = append(vgs, vg)
	}
	return
}

// LogicalVolumes returns a list of LogicalVolumes
func (l *LVM) LogicalVolumes() (lvs LogicalVolumes) {
	for _, lv := range l.lvs {
		lvs = append(lvs, lv)
	}
	return
}

// New return a new and initialized LVM object
func New() (*LVM, error) {
	l := &LVM{
		pvs: make(map[string]*PhysicalVolume),
		vgs: make(map[string]*VolumeGroup),
		lvs: make(map[string]*LogicalVolume),
	}

	if err := l.Read(); err != nil {
		return nil, err
	}

	return l, nil
}

// Read parses the 'pvs' command output and return information about available Physical Volumes,
// Volume Groups and Logical Volumes
func (l *LVM) Read() error {
	var err error
	c := exec.CommandInfo{
		Command: "pvs",
		Args: []string{
			"--noheadings",
			"--nameprefixes",
			"--options=pv_name,pv_size,pv_free,vg_name,vg_free,lv_name,lv_size,lv_tags",
			"--units=b",
		},
	}

	if err = c.Exec(); err != nil || c.Result.Status != 0 {
		return fmt.Errorf("Error whilst querying lvm %v%v", c.Result.Stdout, c.Result.Stderr)
	}

	for _, line := range strings.Split(strings.TrimSuffix(c.Result.Stdout, "\n"), "\n") {
		pvName, vgName, lvName := "", "", ""
		pvSize, pvFree, vgSize, vgFree, lvSize := 0, 0, 0, 0, 0
		vgTags, lvTags := []string{}, []string{}

		for _, field := range strings.Split(line, " ") {
			if len(field) > 0 {
				switch f := strings.Split(field, "="); f[0] {
				case "LVM2_PV_NAME":
					pvName = removeQuotes(f[1])
				case "LVM2_VG_NAME":
					vgName = removeQuotes(f[1])
				case "LVM2_LV_NAME":
					lvName = removeQuotes(f[1])
				case "LVM2_PV_SIZE":
					pvSize, err = getSizeInBytesFromPVS(removeQuotes(f[1]))
					if err != nil {
						return err
					}
				case "LVM2_PV_FREE":
					pvFree, err = getSizeInBytesFromPVS(removeQuotes(f[1]))
					if err != nil {
						return err
					}
				case "LVM2_VG_SIZE":
					vgSize, err = getSizeInBytesFromPVS(removeQuotes(f[1]))
					if err != nil {
						return err
					}
				case "LVM2_VG_FREE":
					vgFree, err = getSizeInBytesFromPVS(removeQuotes(f[1]))
					if err != nil {
						return err
					}
				case "LVM2_LV_SIZE":
					lvSize, err = getSizeInBytesFromPVS(removeQuotes(f[1]))
					if err != nil {
						return err
					}
				case "LVM2_VG_TAGS":
					vgTags = strings.Split(removeQuotes(f[1]), ",")
				case "LVM2_LV_TAGS":
					lvTags = strings.Split(removeQuotes(f[1]), ",")
				default:
					continue
				}
			}
		}

		if _, exists := l.pvs[pvName]; !exists {
			l.pvs[pvName] = &PhysicalVolume{
				Name:        pvName,
				SizeInBytes: pvSize,
				FreeInBytes: pvFree,
			}
		}

		// If we have VG info, create it if it doesn't exist already
		if len(vgName) > 0 {
			if _, exists := l.vgs[vgName]; !exists {
				l.vgs[vgName] = &VolumeGroup{
					Name:        vgName,
					SizeInBytes: vgSize,
					FreeInBytes: vgFree,
					Tags:        vgTags,
				}
			}

			if l.pvs[pvName].VolumeGroup == nil {
				l.pvs[pvName].VolumeGroup = l.vgs[vgName]
			}

			l.vgs[vgName].PhysicalVolumes = append(l.vgs[vgName].PhysicalVolumes, l.pvs[pvName])
		}

		// If we have LV info, create it if it doesn't exist already
		if len(lvName) > 0 {
			if _, exists := l.lvs[lvName]; !exists {
				l.lvs[lvName] = &LogicalVolume{
					Name:        lvName,
					SizeInBytes: lvSize,
					Tags:        lvTags,
					VolumeGroup: l.vgs[vgName],
				}
			}
			l.vgs[vgName].LogicalVolumes = append(l.vgs[vgName].LogicalVolumes, l.lvs[lvName])
		}
	}

	return nil
}

// CreatePhysicalVolume creates a LVM physical volume
func (l *LVM) CreatePhysicalVolume(name string) (*PhysicalVolume, error) {
	c := exec.CommandInfo{
		Command: "pvcreate",
		Args: []string{
			name,
		},
	}

	if err := c.Exec(); err != nil || c.Result.Status != 0 {
		return nil, fmt.Errorf("LVM: Error whilst creating the physical volume %v%v", c.Result.Stdout, c.Result.Stderr)
	}

	if err := l.Read(); err != nil {
		return nil, err
	}

	return l.PhysicalVolume(name)
}

// CreateVolumeGroup creates a LVM volume group
func (l *LVM) CreateVolumeGroup(name string, pvs PhysicalVolumes, tags []string) (*VolumeGroup, error) {
	args := []string{name}
	for _, pv := range pvs {
		args = append(args, pv.Name)
	}

	for _, tag := range tags {
		args = append(args, fmt.Sprintf("--addtag=%s", tag))
	}

	c := exec.CommandInfo{
		Command: "vgcreate",
		Args:    args,
	}

	if err := c.Exec(); err != nil || c.Result.Status != 0 {
		return nil, fmt.Errorf("LVM: Error whilst creating the volume group %v%v", c.Result.Stdout, c.Result.Stderr)
	}

	if err := l.Read(); err != nil {
		return nil, err
	}

	return l.VolumeGroup(name)
}

// CreateLogicalVolume creates a LVM logical volume
func (l *LVM) CreateLogicalVolume(name string, vg *VolumeGroup, sizeInBytes int, tags []string) (*LogicalVolume, error) {
	if sizeInBytes == 0 {
		sizeInBytes = vg.FreeInBytes
	}

	args := []string{
		vg.Name,
		fmt.Sprintf("--size=%db", sizeInBytes),
		fmt.Sprintf("--name=%s", name),
	}

	for _, tag := range tags {
		args = append(args, fmt.Sprintf("--addtag=%s", tag))
	}

	c := exec.CommandInfo{
		Command: "lvcreate",
		Args:    args,
	}

	if err := c.Exec(); err != nil || c.Result.Status != 0 {
		return nil, fmt.Errorf("LVM: Error whilst creating the logical volume %v%v", c.Result.Stdout, c.Result.Stderr)
	}

	if err := l.Read(); err != nil {
		return nil, err
	}

	return l.LogicalVolume(name)
}

// Path returns the device mapper path of a logical volume
func (lv *LogicalVolume) Path() string {
	return fmt.Sprintf("/dev/%s/%s", lv.VolumeGroup.Name, lv.Name)
}

func removeQuotes(str string) string {
	return strings.Replace(str, "'", "", -1)
}

func getSizeInBytesFromPVS(size string) (int, error) {
	return strconv.Atoi(strings.TrimSuffix(removeQuotes(size), "B"))
}
