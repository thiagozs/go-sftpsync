# go-sftpsync - sync folder between sftp/local

`go-sftpsync` is a simple and efficient tool to synchronize files between local and remote directories using SFTP. This tool allows you to upload or download files and supports flexible configurations for connecting to SFTP servers.

## Features

- **Sync Local and Remote Directories**: Easily synchronize files between your local system and an SFTP server.
- **Multiple Profiles**: Manage multiple SFTP configurations using profiles for seamless switching between environments.
- **Flexible Operations**: Perform both **upload** and **download** operations based on your needs.
- **Configurable Authentication**: Supports both password and private key-based authentication.
- **Customizable Paths**: Specify the local and remote base directories to fine-tune which files get synchronized.
- **Error Handling**: Handles errors such as permission denials and connection retries.

## Usage

```bash
go-sftpsync sync [flags]
```

### Flags

- `-c, --config-file string`        Config file path to manage multiple profiles.
- `-h, --help`                      Show help for the `sync` command.
- `-H, --host string`               Specify the SFTP host to connect to.
- `-l, --local-base-path string`    Set the local base path for synchronization.
- `-o, --operation string`          Define the operation to perform: `upload` or `download`.
- `-p, --password string`           SFTP password for authentication.
- `-P, --port string`               Specify the SFTP port (default is 22).
- `-k, --private-key string`        Path to the private key for key-based authentication.
- `-f, --profile string`            Profile name from the configuration file.
- `-r, --remote-base-path string`   Set the remote base path for synchronization.
- `-u, --user string`               SFTP username for authentication.

## Example

To download files from a remote directory to your local machine:

```bash
go-sftpsync sync -o download -H sftp.example.com -u youruser -p yourpassword -l /local/dir -r /remote/dir
```

To upload files from a local directory to an SFTP server:

```bash
go-sftpsync sync -o upload -H sftp.example.com -u youruser -p yourpassword -l /local/dir -r /remote/dir
```

## Configuration File

You can define multiple profiles in a YAML configuration file and switch between them using the `--profile` flag. This simplifies managing multiple SFTP servers or environments.

### Example config.yaml

```yaml
profiles:
  profile1:
    operation: upload
    host: sftp.example.com
    port: "22"
    user: user1
    password: pass1
    private_key: /path/to/private/key1
    local_base_path: /local/path1
    remote_base_path: /remote/path1

  profile2:
    operation: download
    host: sftp.example2.com
    port: "22"
    user: user2
    password: pass2
    private_key: /path/to/private/key2
    local_base_path: /local/path2
    remote_base_path: /remote/path2
```

You can then specify the profile during synchronization:

```bash
go-sftpsync --config-file=$(pwd)/profiles.yaml --profile=backupfiles
```

## License

This project is licensed under the MIT License.

---

© 2024 Thiago Zilli Sarmento. Made with ❤️ and Golang.