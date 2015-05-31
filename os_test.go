package util

func TestCopy(t *testing.T) {
	t.Skip()
	err := Copy(filepath.Join(tmp, "hotfix"), "testdata/hotfix")
	if err != nil {
		panic(err)
	}
}
