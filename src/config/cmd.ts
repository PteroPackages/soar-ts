import { Command } from 'commander';
import { existsSync, readFileSync, writeFileSync } from 'fs';
import { join } from 'path';
import { createInterface } from 'readline';
import log from '../log';
import { getBoolInput } from '../session/response';
import { getConfig, createConfig, getConfigKey, updateConfig } from './funcs';

const infoCmd = new Command('info')
    .description('Gets the Soar configuration setup')
    .addHelpText('before', 'Gets the global Soar configuration setup (or local if specified)')
    .option('--local', 'Gets the local configuration for the workspace', false)
    .option('-h, --hide', 'Hides the API keys from the command output', false)
    .action(async (args: object) => {
        const config = await getConfig(args['local']);
        const appKey = args['hide']
            ? '•'.repeat(config.application.key.length)
            : config.application.key;
        const clientKey = args['hide']
            ? '•'.repeat(config.client.key.length)
            : config.client.key;

        console.log(`
Soar ${args['local'] ? 'Local' : 'Global'} Config
====================
version: ${config.version}

\x1b[4mApplication Details\x1b[0m
url: ${config.application.url || 'Not Set'}
key: ${appKey || 'Not Set'}

\x1b[4mClient Details\x1b[0m
url: ${config.client.url || 'Not Set'}
key: ${clientKey || 'Not Set'}

\x1b[4mLogging\x1b[0m
show debug logs:     ${config.logs.showDebug}
show http logs:      ${config.logs.showHttp}
show websocket logs: ${config.logs.showWebsocket}
use colour:          ${config.logs.useColour}

\x1b[4mHTTP\x1b[0m
retry on ratelimit:      ${config.http.retryRatelimit}
save http requests:      ${config.http.saveRequests}
send full response body: ${config.http.sendFullBody}

\x1b[4mCore\x1b[0m
ignore warnings:       ${config.core.ignoreWarnings}
save error logs:       ${config.core.saveErrorLogs}
stop at system errors: ${config.core.stopAtSysError}
`);
    });

const setupCmd = new Command('setup')
    .description('Setup a new Soar configuration')
    .addHelpText('before', 'Setup a new global or local Soar configuration')
    .option('--local', 'Setup a local configuration for the workspace', false)
    .option('--link [file]', 'Links the new config with another local config or the global config if not provided')
    .option('-f, --force', 'Skips all confirmation prompts', false)
    .action(async (args: object) => {
        const local: boolean = args['local'];
        const force: boolean = args['force'];
        let link: string | boolean = args['link'];
        let mainfp = join(process.env.SOAR_PATH || '', 'config.yml');

        if (link) {
            if (typeof link === 'string') {
                if (!existsSync(link)) log.error(
                    'Not Found Error',
                    'the local config file path could not be resolved',
                    true
                );
                mainfp = link;
            } else {
                if (!existsSync(process.env.SOAR_PATH)) log.error('MISSING_CONFIG', null, true);
            }
        }

        if (local) {
            if (existsSync(join(process.cwd(), '.soar-local.yml'))) {
                log.info('existing local config file found');
                if (!force) {
                    const reader = createInterface(
                        process.stdin,
                        process.stdout
                    );

                    const res = await getBoolInput(reader, 'do you want to overwrite this file? (y/n)');
                    if (!res) return;
                } else {
                    log.info('overwrite mode forced for local config');
                }
            }

            if (!existsSync(mainfp)) log.error('MISSING_ENV', null, true);
            try {
                const linkData = readFileSync(mainfp, { encoding: 'utf-8' });
                writeFileSync(
                    join(process.cwd(), '.soar-local.yml'),
                    linkData, { encoding: 'utf-8' }
                );

                log.success([
                    'setup a new local config at:',
                    join(process.cwd(), '.soar-local.yml')
                ]);
            } catch (err) {
                log.fromError(err, true);
            }
        } else {
            await createConfig(link ? mainfp : null);
        }
    });

const setCmd = new Command('set')
    .description('Sets an option in the Soar configuration')
    .addHelpText('before', 'Sets an option in the Soar configuration')
    .argument('<key>', 'The config key to set')
    .argument('<value>', 'The value to set the key to')
    .option('--local', 'Updates the local Soar configuration', false)
    .action(async (key: string, value: string, args: object) => {
        const config = await getConfig(args['local']);
        const option = getConfigKey(config, key);

        if (!option.length) {
            log.error('Invalid Arguments', `invalid config key '${key}'`);
            log.notice("you can view all config keys with the '%bsoar config info%R' command");
            return;
        }

        updateConfig(config, option, value, args['local']);
        log.success('updated config');
    });

export default [
    infoCmd,
    setupCmd,
    setCmd
]
