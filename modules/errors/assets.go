// Code generated by go-bindata.
// sources:
// assets/errors.csv
// DO NOT EDIT!

package errors

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _assetsErrorsCsv = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x84\x90\x31\xd3\xd3\x30\x0c\x86\xf7\xfe\x0a\x6d\x2c\xb9\xc4\xbd\x36\x0b\x1b\xd7\x63\x61\x00\xee\xca\xc4\xe6\xd8\x2f\x8d\x69\x2a\x19\x59\x49\x5b\x7e\x3d\xe7\xb4\x70\xc7\xf0\x5d\x47\xcb\x8f\xf4\x3e\x52\xef\x9c\x73\x4d\xef\x5c\x33\x9a\xe5\xf7\x5d\x87\x9b\xbf\xe4\x09\x6d\x90\x4b\x77\x4d\xe7\xd4\x41\x55\xb4\x74\x0f\xf0\x08\x5d\xa0\xf4\xb1\xd6\x9a\x0f\x4c\x33\x47\x18\xf4\x92\x18\x91\x56\x92\x46\x9f\x33\xea\x53\x98\x6c\x04\x3d\x5a\xda\x4d\xef\x76\x6e\xdb\xf4\x6e\xf7\x3a\xa9\x82\x07\xcf\x2c\x46\x41\x98\x11\x8c\x4c\xe8\xbb\xc8\x19\xc8\xd0\xe6\x20\xf3\x14\xe9\xad\x6f\xf2\x46\xa7\xb4\x80\xc9\xc7\xa8\x28\xa5\xa5\xc3\x88\x70\x26\x1b\xbd\xd1\xef\x7f\x58\x2a\xa4\x33\x73\xe2\x13\x79\x8e\xe4\x43\x40\x29\x69\x98\x40\xc3\x7d\x35\x57\x14\xa3\xf2\xd4\xdf\x3b\xe7\xb6\xcd\xfe\xf5\xa1\x1e\xe0\x53\x3f\x7b\x2d\xa0\x41\xe2\x9d\x12\x9b\x10\x6e\x19\xc1\x10\xe9\xd3\xf1\xcb\x67\x92\xe1\x27\x82\x35\xdf\xd6\xac\x5f\x73\x8d\x8b\x82\xc2\xef\xea\x62\xaa\x28\x59\x38\xd6\xdd\xaa\x4d\x31\x9d\x83\xcd\x0a\x62\x20\x62\xad\x17\x99\x16\xd0\x5d\x66\xfd\x3b\xa0\xa5\xaf\x13\x7c\xa9\x03\x97\x84\xeb\xda\xb9\xc6\xcb\x8f\xff\xb9\xcd\x9f\x00\x00\x00\xff\xff\x53\xba\x35\x5f\xfa\x01\x00\x00")

func assetsErrorsCsvBytes() ([]byte, error) {
	return bindataRead(
		_assetsErrorsCsv,
		"assets/errors.csv",
	)
}

func assetsErrorsCsv() (*asset, error) {
	bytes, err := assetsErrorsCsvBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/errors.csv", size: 506, mode: os.FileMode(420), modTime: time.Unix(1489478586, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"assets/errors.csv": assetsErrorsCsv,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"assets": &bintree{nil, map[string]*bintree{
		"errors.csv": &bintree{assetsErrorsCsv, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
