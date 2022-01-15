import { FlagOptions } from '../structs';

export function parseUserGroup(args: object): FlagOptions {
    const type = (args['text'] && 'text') || (args['yaml'] && 'yaml') || 'json';
    let file = '';

    if (args['output']) {
        if (typeof args['output'] === 'boolean') file = `soar_log_${Date.now()}`;
        else file = args['output'];
    }
    if (file.length && !file.endsWith('.'+ type)) file += '.'+ type;

    return {
        silent: args['silent'],
        prompt: args['prompt'],
        writeFile: file,
        responseType: type
    } as FlagOptions;
}

export function buildUser(args: object): string {
    let base = '/api/application/users';
    if (args['id']) return `${base}/${args['id']}`;
    if (args['email']) return `${base}?filter[email]=${args['email']}`;
    if (args['uuid']) return `${base}?filter[uuid]=${args['uuid']}`;
    if (args['username']) return `${base}?filter[username]=${args['username']}`;
    if (args['external']) return `${base}?filter[external_id]=${args['external']}`;
    return base;
}
