# aipchecker

- Multithreaded unpacking of zip-archived Android projects
- Multithreaded checking if apps require Internet permission

## Usage

```
aipchecker [flags] [path/to/Android/projects/folder]
```

### Flags

The -unzip flag indicates if projects need to be unzipped first

### Examples

```
aipchecker -unzip "~/documents/Android Projects"
aipchecker "~/documents/Android Projects"
```

## License

MIT