package bundle

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/afs"
	"github.com/mangalorg/mangal/provider/info"
	"github.com/mangalorg/mangal/provider/lua"
)

func Loaders(dir string) ([]libmangal.ProviderLoader, error) {
	dirEntries, err := afs.Afero.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var bundleLoaders []libmangal.ProviderLoader
	for _, dirEntry := range dirEntries {
		// skip non directories
		if !dirEntry.IsDir() {
			continue
		}

		dirEntryPath := filepath.Join(dir, dirEntry.Name())

		isProvider, err := afs.Afero.Exists(filepath.Join(dirEntryPath, info.Filename))
		if err != nil {
			return nil, err
		}

		if !isProvider {
			continue
		}

		loaders, err := getLoaders(dirEntryPath)
		if err != nil {
			return nil, err
		}

		bundleLoaders = append(bundleLoaders, loaders...)
	}

	return bundleLoaders, nil
}

func singletone(loader libmangal.ProviderLoader) []libmangal.ProviderLoader {
	return []libmangal.ProviderLoader{loader}
}

func getLoaders(dir string) ([]libmangal.ProviderLoader, error) {
	infoFile, err := afs.Afero.OpenFile(
		filepath.Join(dir, info.Filename),
		os.O_RDONLY,
		0755,
	)

	if err != nil {
		return nil, err
	}

	defer infoFile.Close()

	providerInfo, err := info.Parse(infoFile)

	if err != nil {
		return nil, err
	}

	switch providerInfo.Type {
	case info.TypeLua:
		loader, err := lua.NewLoader(providerInfo.Provider, dir)
		if err != nil {
			return nil, err
		}

		return singletone(loader), nil
	case info.TypeBundle:
		return Loaders(dir)
	default:
		return nil, fmt.Errorf("unkown provider type: %#v", providerInfo.Type)
	}
}
