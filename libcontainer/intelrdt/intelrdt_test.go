package intelrdt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntelRdtSetL3CacheSchema(t *testing.T) {
	helper := NewIntelRdtTestUtil(t)

	const (
		l3CacheSchemaBefore = "L3:0=f;1=f0"
		l3CacheSchemeAfter  = "L3:0=f0;1=f"
	)

	helper.writeFileContents(map[string]string{
		"schemata": l3CacheSchemaBefore + "\n",
	})

	helper.config.IntelRdt.L3CacheSchema = l3CacheSchemeAfter
	intelrdt := newManager(helper.config, "", helper.IntelRdtPath)
	if err := intelrdt.Set(helper.config); err != nil {
		t.Fatal(err)
	}

	tmpStrings, err := getIntelRdtParamString(helper.IntelRdtPath, "schemata")
	if err != nil {
		t.Fatalf("Failed to parse file 'schemata' - %s", err)
	}
	values := strings.Split(tmpStrings, "\n")
	value := values[0]

	if value != l3CacheSchemeAfter {
		t.Fatal("Got the wrong value, set 'schemata' failed.")
	}
}

func TestIntelRdtSetMemBwSchema(t *testing.T) {
	helper := NewIntelRdtTestUtil(t)

	const (
		memBwSchemaBefore = "MB:0=20;1=70"
		memBwSchemeAfter  = "MB:0=70;1=20"
	)

	helper.writeFileContents(map[string]string{
		"schemata": memBwSchemaBefore + "\n",
	})

	helper.config.IntelRdt.MemBwSchema = memBwSchemeAfter
	intelrdt := newManager(helper.config, "", helper.IntelRdtPath)
	if err := intelrdt.Set(helper.config); err != nil {
		t.Fatal(err)
	}

	tmpStrings, err := getIntelRdtParamString(helper.IntelRdtPath, "schemata")
	if err != nil {
		t.Fatalf("Failed to parse file 'schemata' - %s", err)
	}
	values := strings.Split(tmpStrings, "\n")
	value := values[0]

	if value != memBwSchemeAfter {
		t.Fatal("Got the wrong value, set 'schemata' failed.")
	}
}

func TestIntelRdtSetCombinedSchema(t *testing.T) {
	helper := NewIntelRdtTestUtil(t)

	// Test filtering out the "MB:"" line in l3CacheSchema.

	const (
		schemaBefore        = "MB:0=20;1=70"
		memBwSchema         = "MB:0=70;1=20"
		l3CacheSchema       = "MB:0=80;1=10\nL3:0=f0;1=f"
		combinedSchemaAfter = "L3:0=f0;1=f\nMB:0=70;1=20"
	)

	helper.writeFileContents(map[string]string{
		"schemata": schemaBefore + "\n",
	})

	helper.config.IntelRdt.MemBwSchema = memBwSchema
	helper.config.IntelRdt.L3CacheSchema = l3CacheSchema
	intelrdt := newManager(helper.config, "", helper.IntelRdtPath)
	if err := intelrdt.Set(helper.config); err != nil {
		t.Fatal(err)
	}

	tmpStrings, err := getIntelRdtParamString(helper.IntelRdtPath, "schemata")
	if err != nil {
		t.Fatalf("Failed to parse file 'schemata' - %s", err)
	}

	readValues := strings.Split(tmpStrings, "\n")
	expectedValues := strings.Split(combinedSchemaAfter, "\n")

	if readValues[0] != expectedValues[0] {
		t.Fatal("Got the wrong value for L3 cache, set 'schemata' failed.")
	}

	if readValues[1] != expectedValues[1] {
		t.Fatal("Got the wrong value for MemBW, set 'schemata' failed.")
	}
}

func TestIntelRdtSetMemBwScSchema(t *testing.T) {
	helper := NewIntelRdtTestUtil(t)

	const (
		memBwScSchemaBefore = "MB:0=5000;1=7000"
		memBwScSchemeAfter  = "MB:0=9000;1=4000"
	)

	helper.writeFileContents(map[string]string{
		"schemata": memBwScSchemaBefore + "\n",
	})

	helper.config.IntelRdt.MemBwSchema = memBwScSchemeAfter
	intelrdt := newManager(helper.config, "", helper.IntelRdtPath)
	if err := intelrdt.Set(helper.config); err != nil {
		t.Fatal(err)
	}

	tmpStrings, err := getIntelRdtParamString(helper.IntelRdtPath, "schemata")
	if err != nil {
		t.Fatalf("Failed to parse file 'schemata' - %s", err)
	}
	values := strings.Split(tmpStrings, "\n")
	value := values[0]

	if value != memBwScSchemeAfter {
		t.Fatal("Got the wrong value, set 'schemata' failed.")
	}
}

func TestSchemataLineComparison(t *testing.T) {
	const (
		appliedMemBwSchema                             = "MB:0=5000;1=7000"
		existingMemBwSchema                            = "MB:0=5000;1=7000"
		appliedMemBwSchemaInDifferentOrder             = "MB:1=7000;0=5000;"
		appliedL3CacheSchema                           = "L3:0=f0;1=f"
		existingL3CacheSchema                          = "L3:0=f0;1=f"
		appliedL3CacheSchemaWithSpaces                 = "L3: 0= f0; 1 =   f"
		existingL3CacheSchemaWithSpaces                = "L3:  0 = f0 ;  1 =f"
		appliedL3CacheSchemaWithSpacesAndLeadingZeros  = "L3: 0= f0; 1 =   000f"
		existingL3CacheSchemaWithSpacesAndLeadingZeros = "L3:  0 = 00f0 ;  1 =f"
		existingL3CacheSchemaInDifferentOrder          = "L3:  1=f; 0 = 00f0 ;"
		brokenL3CacheSchema                            = "L3:  1f; 0 = 00f0 ;"
		differentL3CacheSchema                         = "L3: 0=f; 1 =f"
		appliedMultiLineSchema                         = "\nMB:0=70;1 = 20\nL3:1=   f ;0 = 00f0;"
		existingMultiLineSchema                        = "L3:1=000f;0= f0\nMB: 0=  70;1=20\n\n"
		existingDifferentMultiLineSchema               = "L3:0   =f0;1=000f\nMB: 0=  71;1=20\n\n"
	)

	err := checkSchemataMatch(appliedMemBwSchema, existingMemBwSchema)
	if err != nil {
		t.Fatal("MemBw schematas should have matched.")
	}

	err = checkSchemataMatch(appliedMemBwSchemaInDifferentOrder, existingMemBwSchema)
	if err != nil {
		t.Fatal("MemBw schematas in different order should have matched.")
	}

	err = checkSchemataMatch(appliedL3CacheSchema, existingL3CacheSchema)
	if err != nil {
		t.Fatal("L3Cache schematas should have matched.")
	}

	err = checkSchemataMatch(appliedL3CacheSchemaWithSpaces, existingL3CacheSchemaWithSpaces)
	if err != nil {
		t.Fatal("L3Cache schematas with spaces should have matched.")
	}

	err = checkSchemataMatch(appliedL3CacheSchemaWithSpacesAndLeadingZeros, existingL3CacheSchemaWithSpacesAndLeadingZeros)
	if err != nil {
		t.Fatal("L3Cache schematas with spaces and leading zeros should have matched.")
	}

	err = checkSchemataMatch(appliedL3CacheSchema, existingL3CacheSchemaWithSpacesAndLeadingZeros)
	if err != nil {
		t.Fatal("L3Cache schematas with spaces and leading zeros should have matched to the original string.")
	}

	err = checkSchemataMatch(appliedL3CacheSchema, existingL3CacheSchemaInDifferentOrder)
	if err != nil {
		t.Fatal("L3Cache schematas in different order should have matched to the original string.")
	}

	err = checkSchemataMatch(existingL3CacheSchema, brokenL3CacheSchema)
	if err == nil {
		t.Fatal("L3Cache schemata should have failed to parse")
	}

	err = checkSchemataMatch(appliedL3CacheSchema, differentL3CacheSchema)
	if err == nil {
		t.Fatal("L3Cache schemata should not have matched")
	}

	err = checkSchemataMatch(appliedMultiLineSchema, existingMultiLineSchema)
	if err != nil {
		t.Fatal("Multi-line schematas should have matched.")
	}

	err = checkSchemataMatch(appliedMultiLineSchema, existingDifferentMultiLineSchema)
	if err == nil {
		t.Fatal("Multi-line schematas should not have matched.")
	}
}

func TestApply(t *testing.T) {
	helper := NewIntelRdtTestUtil(t)

	const closID = "test-clos"

	helper.config.IntelRdt.ClosID = closID
	intelrdt := newManager(helper.config, "", helper.IntelRdtPath)
	if err := intelrdt.Apply(1234); err == nil {
		t.Fatal("unexpected success when applying pid")
	}
	if _, err := os.Stat(filepath.Join(helper.IntelRdtPath, closID)); err == nil {
		t.Fatal("closid dir should not exist")
	}

	// Dir should be created if some schema has been specified
	intelrdt.config.IntelRdt.L3CacheSchema = "L3:0=f"
	if err := intelrdt.Apply(1235); err != nil {
		t.Fatalf("Apply() failed: %v", err)
	}

	pids, err := getIntelRdtParamString(intelrdt.GetPath(), "tasks")
	if err != nil {
		t.Fatalf("failed to read tasks file: %v", err)
	}
	if pids != "1235" {
		t.Fatalf("unexpected tasks file, expected '1235', got %q", pids)
	}
}

func TestIntelRdtSetWithClosId(t *testing.T) {
	helper := NewIntelRdtTestUtil(t)

	const (
		closID                 = "test-clos"
		l3CacheSchemaBefore    = "L3:0=f;1=f0"
		l3CacheSchemeDifferent = "L3:0=f0;1=f"
		l3CacheSchemeSame      = "L3:1=f0;0=f"
	)

	helper.config.IntelRdt.ClosID = closID
	intelrdt := newManager(helper.config, "", helper.IntelRdtPath)

	// Create an "existing" schmata file to match against.
	helper.writeFileContents(map[string]string{
		"schemata": l3CacheSchemaBefore + "\n",
	})

	helper.config.IntelRdt.L3CacheSchema = l3CacheSchemeDifferent
	if err := intelrdt.Set(helper.config); err == nil {
		t.Fatalf("expected the file comparison to fail")
	}

	helper.config.IntelRdt.L3CacheSchema = l3CacheSchemeSame
	if err := intelrdt.Set(helper.config); err != nil {
		t.Fatalf("expected the file comparison to succeed: %s", err)
	}
}
