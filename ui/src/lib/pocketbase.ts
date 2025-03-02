import PocketBase from 'pocketbase';

const baseUrl: string = import.meta.env.MODE == 'development' ? 'http://127.0.0.1:8090' : '/';
const pb = new PocketBase(baseUrl);

// export default { pb: pb };
export default pb;
