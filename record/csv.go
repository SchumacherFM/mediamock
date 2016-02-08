package record

import (
	"compress/gzip"
	"encoding/csv"
	"io"
	"net/http"
	"os"

	"github.com/SchumacherFM/mediamock/common"
)

func GetCSVContent(csvFile string) [][]string {
	var rawRC io.ReadCloser
	if common.IsHTTP(csvFile) {
		resp, err := http.Get(csvFile)
		if err != nil {
			common.UsageAndExit("Failed to download %q with error: %s", csvFile, err)
		}
		if resp.StatusCode != http.StatusOK {
			common.UsageAndExit("Server return non-200 status code: %s\nFailed to download %s", resp.Status, csvFile)
		}
		rawRC = resp.Body
	} else {
		var err error
		rawRC, err = os.Open(csvFile)
		if err != nil {
			common.UsageAndExit("Failed to open %q with error:%s", csvFile, err)
		}
	}
	defer func() {
		if err := rawRC.Close(); err != nil {
			common.UsageAndExit("Failed to close URL/file %s with error: %s", csvFile, err)
		}
	}()

	rz, err := gzip.NewReader(rawRC)
	if err != nil {
		common.UsageAndExit("Failed to create a GZIP reader from file %s with error: %s", csvFile, err)
	}
	defer func() {
		if err := rz.Close(); err != nil {
			common.UsageAndExit("Failed to close file %q with error: %s", csvFile, err)
		}
	}()

	rc := csv.NewReader(rz)
	rc.Comma = ([]rune(CSVSep))[0]

	records, err := rc.ReadAll()
	if err != nil {
		common.UsageAndExit("Failed to read CSV file %q with error: %s", csvFile, err)
	}

	return records
}
