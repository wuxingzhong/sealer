// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package buildimage

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"helm.sh/helm/v3/pkg/chartutil"

	"github.com/sealerio/sealer/common"
	osi "github.com/sealerio/sealer/utils/os"
	strUtils "github.com/sealerio/sealer/utils/strings"
)

func getKubeVersion(rootfs string) string {
	chartsPath := filepath.Join(rootfs, "charts")
	if !osi.IsFileExist(chartsPath) {
		return ""
	}
	return readCharts(chartsPath)
}

func readCharts(chartsPath string) string {
	var kv string
	err := filepath.Walk(chartsPath, func(path string, f fs.FileInfo, err error) error {
		if kv != "" {
			return nil
		}
		if f.IsDir() || f.Name() != "Chart.yaml" {
			return nil
		}
		meta, walkErr := chartutil.LoadChartfile(path)
		if walkErr != nil {
			return walkErr
		}
		if meta.KubeVersion != "" {
			kv = meta.KubeVersion
		}
		return nil
	})

	if err != nil {
		return ""
	}
	return kv
}

func FormatImages(images []string) (res []string) {
	for _, img := range images {
		tmpImg := strings.TrimSpace(img)
		tmpImg = strings.Trim(tmpImg, `'"`)
		tmpImg = strings.TrimSpace(tmpImg)
		if strings.HasPrefix(tmpImg, "#") || tmpImg == "" {
			continue
		}
		res = append(res, tmpImg)
	}

	res = strUtils.RemoveDuplicate(res)
	return
}

//func trimQuotes(s string) string {
//	if len(s) >= 2 {
//		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
//			return s[1 : len(s)-1]
//		}
//	}
//	return s
//}

func marshalJSONToFile(file string, obj interface{}) error {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	if err = os.WriteFile(file, data, common.FileMode0644); err != nil {
		return err
	}
	return nil
}
