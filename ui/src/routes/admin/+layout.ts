import pb from '$lib/pocketbase';

export const load = async () => {
	try {
		return {
			user: pb.authStore.model
		};
	} catch {
		return {
			user: null
		};
	}
};
