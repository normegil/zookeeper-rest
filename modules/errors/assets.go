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

var _assetsErrorsCsv = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x90\xbd\x6e\x1c\x31\x0c\x84\x7b\x3f\x05\xbb\x34\xc2\xae\xfc\x73\x4d\x3a\xe3\x90\x26\x40\x7e\x00\xa7\x4a\xa7\x93\xc6\x5e\xe5\x76\x49\x85\xa2\xd6\xbe\x3c\x7d\xa0\xbd\x4b\x00\x17\x81\x53\x09\x22\x3f\x72\x66\xb8\xf3\xde\x7b\xb7\xf3\xde\x4d\x66\xe5\xfd\x38\xe2\x25\x2c\x65\xc6\x10\x65\x19\x9f\xf3\x31\x8f\x50\x15\xad\xe3\x19\x7c\x80\xae\x50\xfa\xd0\x6b\xee\x9e\xa9\x71\x82\x41\x97\xcc\x48\xb4\x91\x34\x85\x52\xd0\xbf\xc2\x64\x13\xe8\x3c\x32\x5c\xed\xfc\xad\xbf\x76\x3b\x7f\xfb\xb6\x52\x07\xf7\x81\x59\x8c\xa2\x30\x23\x1a\x99\xd0\x77\x91\x23\x50\xa0\x6e\x2f\x6d\x4e\xf4\xaf\x36\x05\xa3\xa7\xbc\x82\x29\xa4\xa4\xa8\x75\xa0\xfd\x84\x78\x24\x9b\x82\xd1\xaf\xbf\x58\xae\xa4\x8d\x39\xf3\x13\x05\x4e\x14\x62\x44\xad\xf9\x30\x83\x0e\xa7\xcd\xb9\xa2\x1a\xd5\x8b\xfd\x3b\xef\xfd\xb5\xbb\x7b\xfb\x50\x67\xf0\x62\xbf\x04\xad\xa0\x83\xa4\x13\x65\x36\x21\xbc\x14\x44\x43\xa2\x8f\x0f\x5f\x3e\x93\x1c\x7e\x20\x9a\xfb\xb6\x69\xfd\x6c\x5d\x2e\x09\x2a\xbf\xeb\xc1\x54\x51\x8b\x70\xea\xd9\xba\x9b\x6a\xda\xa2\x35\x05\x31\x90\xb0\xd5\xab\xcc\x2b\xe8\x24\x4d\xff\x2c\x18\xe8\xeb\x8c\x50\xfb\xc2\x35\xe3\x79\x9b\xdc\xe4\xe5\xf1\x35\xb7\x05\xba\xf9\xdf\x40\x37\xee\x53\xae\xb5\x9f\xaa\xcf\x67\x45\xa2\x35\xcc\x0d\xee\xfe\xfc\xf6\x63\x2e\x17\xe2\x51\x65\x79\xa5\x45\xa2\x97\x7e\x99\x43\x44\x1a\xae\x7e\x07\x00\x00\xff\xff\x7b\xf3\x49\x6e\x76\x02\x00\x00")

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

	info := bindataFileInfo{name: "assets/errors.csv", size: 630, mode: os.FileMode(420), modTime: time.Unix(1489503474, 0)}
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

