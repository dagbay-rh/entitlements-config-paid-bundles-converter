# entitlements-config-paid-bundles-converter

A Go program that converts entitlements config bundles to paid/not paid bundles by filtering out paid SKUs from the regular SKUs list.
This is meant to run locally as a one time thing with config files from [entitlements-config](https://github.com/RedHatInsights/entitlements-config) 
while we move the isEval logic from feature service to entitlements service.

The intention here is to manually update existing config files by adding the appropriate `paid_skus` 
field to each bundle. Then run this program to remove any of the skus in `paid_skus` from `skus`. Any bundle without `paid_skus` will be left unaltered, and the other fields on each bundle will be preserved.

## Usage

### Running the Program

```bash
go run main.go [OPTIONS] <filepath>
```

**Arguments:**
- `<filepath>` - Path to the bundles ConfigMap template file

**Options:**
- `-v`, `--verbose` - Enable verbose output (prints loaded bundles)
- `-h`, `--help` - Show help message

**Examples:**

Run with a bundles file:
```bash
go run main.go /home/janedoe/entitlements-config/configs/stage/bundles.yml
```

The program will create a new file with the `.new` extension containing the filtered bundles.

Example output (continuing from above):
```bash
Successfully wrote filtered bundles to: /home/janedoe/entitlements-config/configs/stage/bundles.yml.new
```

Validate manually, and then you can replace the existing config with the new one and commit/push.

### Example input/ouput
Input file:
```yaml
kind: Template
apiVersion: v1
objects:
    - apiVersion: v1
      kind: ConfigMap
      metadata:
        annotations:
            qontract.recycle: "true"
        creationTimestamp: null
        name: entitlements-config
      data:
        bundles.yml: |
            - name: foo
              skus:
                - RH0001
                - RH0002
                - RH0003
              paid_skus:
                - RH0002
                - RH0003

            - name: bar
              skus:
                - RH0004
                - RH0005
```

Output file:
```yaml
kind: Template
apiVersion: v1
objects:
    - apiVersion: v1
      kind: ConfigMap
      metadata:
        annotations:
            qontract.recycle: "true"
        creationTimestamp: null
        name: entitlements-config
      data:
        bundles.yml: |
            - name: foo
              skus:
                - RH0001
              paid_skus:
                - RH0002
                - RH0003

            - name: bar
              skus:
                - RH0004
                - RH0005
```


### Tests

Run tests with verbose output:
```bash
go test -v
```

## Disclaimer

This program was written with the assistance of Claude (Sonnet 4.5).