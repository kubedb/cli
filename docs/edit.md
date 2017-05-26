# kubedb edit

## Example

##### Help for edit command

```bash
$ kubedb edit --help

Edit a resource from the default editor.

The edit command allows you to directly edit any API resource you can retrieve via the command line tools. It will open
the editor defined by your KUBEDB _EDITOR, or EDITOR environment variables, or fall back to 'nano'

Examples:
  # Edit the elastic named 'elasticsearch-demo':
  kubedb edit es/elasticsearch-demo

  # Use an alternative editor
  KUBEDB_EDITOR="nano" kubedb edit es/elasticsearch-demo

Options:
  -n, --namespace='default': Edit object(s) in this namespace.
  -o, --output='yaml': Output format. One of: yaml|json.

Usage:
  kubedb edit (RESOURCE/NAME) [flags] [options]

Use "kubedb edit options" for a list of global command-line options (applies to all commands).
```
