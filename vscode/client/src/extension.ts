import * as path from 'path';
import { workspace, ExtensionContext } from 'vscode';

import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: ExtensionContext) {
  // The debug options for the server
  // --inspect=6009: runs the server in Node's Inspector mode so VS Code can attach to the server for debugging
  const debugOptions = { execArgv: ['--nolazy', '--inspect=6009'] };

  // If the extension is launched in debug mode then the debug server options are used
  // Otherwise the run options are used
  const serverOptions: ServerOptions = {
    run: {
      command: 'risor-lsp', // FIXME: Need to install this binary automatically
      args: [],
      options: {},
    },
    debug: {
      command: 'risor-lsp',
      args: [],
      options: {},
    },
  };

  // Options to control the language client
  const clientOptions: LanguageClientOptions = {
    // Register the server for Risor files
    documentSelector: [
      { scheme: 'file', language: 'risor' },
      { scheme: 'file', pattern: '**/*.risor' },
      { scheme: 'file', pattern: '**/*.rsr' }
    ],
    synchronize: {
      // Notify the server about file changes to Risor files
      fileEvents: workspace.createFileSystemWatcher('**/*.{risor,rsr}'),
    },
  };

  // Create the language client and start the client.
  client = new LanguageClient(
    'risorLanguageServer',
    'Risor Language Server',
    serverOptions,
    clientOptions
  );

  // Start the client. This will also launch the server
  client.start();
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
