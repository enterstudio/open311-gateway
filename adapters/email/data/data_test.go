package data

import (
	"fmt"
	"testing"

	"Gateway311/adapters/email/logs"
	"Gateway311/adapters/email/structs"

	"github.com/davecgh/go-spew/spew"
)

var Debug = true

func init() {
	logs.Init(Debug)

	fmt.Println("Reading config...")
	if err := Init("config.json"); err != nil {
		fmt.Printf("Init() failed: %s", err)
	}
}

type testResultS struct {
	input string
	isOK  bool
}

func isOK(e error) bool {
	if e == nil {
		return false
	}
	return true
}

func TestServices(t *testing.T) {
	fmt.Printf("\n\n\n\n============================= [TestServices] =============================\n\n")

	var test1 = [3]testResultS{
		{"Cupertino", true},
		{"San Jose", false},
		{"Sunnyvale", true},
	}

	for _, tt := range test1 {
		svcs, err := ServicesArea(tt.input)

		switch {
		case tt.isOK && err == nil:
			fmt.Printf("svcs for %q:\n%s", tt.input, spew.Sdump(*svcs))
		case tt.isOK && err != nil:
			t.Errorf("ServicesArea() failed for: %q", tt.input)
		case !tt.isOK && err == nil:
			t.Errorf("ServicesArea() should have failed for: %q", tt.input)
		}
	}

	fmt.Printf("----------------------------- [TestServicesAll] -----------------------------\n\n")
	svcs, err := ServicesAll()
	if err != nil {
		t.Errorf("ServicesArea() failed.")
	} else {
		fmt.Printf("svcs:\n%s", spew.Sdump(*svcs))
	}

}

func TestAdapter(t *testing.T) {
	fmt.Printf("\n\n\n\n============================= [TestAdapter] =============================\n\n")

	fmt.Printf("----------------------------- [Adapter] -----------------------------\n\n")
	if name, atype, address := Adapter(); name != "EM1" || atype != "Email" || address != "" {
		t.Errorf("Adapter() failed - name: %q  atype: %q  address: %q", name, atype, address)
	} else {
		fmt.Println("OK!")
	}

	fmt.Printf("----------------------------- [AdapterName] -----------------------------\n\n")
	if name := AdapterName(); name != "EM1" {
		t.Errorf("AdapterName() failed - name: %q", name)
	} else {
		fmt.Println("OK!")
	}

	fmt.Printf("\n\n----------------------------- [MIDProvider] -----------------------------\n\n")
	var midTests = []struct {
		n    structs.ServiceID // input
		isOK bool              // expected result
	}{
		{structs.ServiceID{"EM1", "CU", 1, 1}, true},
		{structs.ServiceID{"EM1", "CU", 2, 9999}, true},
		{structs.ServiceID{"EM1", "CU", 3, 999999}, false},
		{structs.ServiceID{"EM1", "SUN", 1, 1}, true},
		{structs.ServiceID{"EM1", "SJ", 1, 1}, false},
		{structs.ServiceID{"EM1", "XXXXXXXXX", 1, 1}, false},
	}

	for _, tt := range midTests {
		prov, err := MIDProvider(tt.n)
		switch {
		case tt.isOK && err == nil:
			fmt.Printf("\nsvcs for %q:\n%s", tt.n.MID(), prov.String())
		case tt.isOK && err != nil:
			t.Errorf("ServicesArea() failed for: %q", tt.n.MID())
		case !tt.isOK && err == nil:
			t.Errorf("ServicesArea() should have failed for: %q", tt.n)
		}
	}

	fmt.Printf("\n\n----------------------------- [RouteProvider] -----------------------------\n\n")
	var routeTests = []struct {
		n    structs.NRoute // input
		isOK bool           // expected result
	}{
		{structs.NRoute{"EM1", "CU", 1}, true},
		{structs.NRoute{"EM1", "CU", 2}, true},
		{structs.NRoute{"EM1", "CU", 3}, false},
		{structs.NRoute{"EM1", "SUN", 1}, true},
		{structs.NRoute{"EM1", "SJ", 1}, false},
		{structs.NRoute{"EM1", "XXXXXXXXX", 1}, false},
	}

	for _, tt := range routeTests {
		prov, err := RouteProvider(tt.n)
		switch {
		case tt.isOK && err == nil:
			fmt.Printf("\nsvcs for %q:\n%s", tt.n.SString(), prov.String())
		case tt.isOK && err != nil:
			t.Errorf("ServicesArea() failed for: %q", tt.n.SString())
		case !tt.isOK && err == nil:
			t.Errorf("ServicesArea() should have failed for: %q", tt.n)
		}
	}
}