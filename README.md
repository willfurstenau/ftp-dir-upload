# ftp-dir-upload
Github action to upload static directory of files to an FTP server
## Inputs

### `host`

**Required** The FTP host.

### `port`

**Default: 21** FTP Port

### `user`

**Required** FTP username

### `password`

**Required** FTP password

### `local-dir`

**Required** Local directory to upload (relative to `$GITHUB_WORKSPACE`)

### `remote-dir`

**Required** Remote directory on the FTP server (will be created if missing)

## Example usage
```yaml
- name: FTP dir upload
  uses: willfurstenau/ftp-dir-upload@v1
  with:
    host: ftp.myhost.com
    user: ${{ secrets.FTP_USER }}
    password: ${{ secrets.FTP_PASS }}
    local-dir: ./dist
    remote-dir: /public_html
```