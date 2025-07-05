import { goto } from '$app/navigation';
import pb from '$lib/pocketbase';

export const logOutUser = async () => {
	pb.authStore.clear();
	await goto('/login');
};

export const isLoggedIn = () => {
	return pb.authStore.isValid;
};

export const isAdmin = () => {
	return isLoggedIn() && pb.authStore?.model?.role == 'admin';
};

export const loggedInUserID = () => {
	return isLoggedIn() && pb.authStore.model?.id;
};

export const objectMap = (obj, fn) => {
	obj ||= {};

	const newObject = {};
	Object.keys(obj).forEach((key) => {
		newObject[key] = fn(obj[key]);
	});
	return newObject;
};

export const isComponent = (value: any) => {
	return typeof value === 'function' && /^class\s/.test(Function.prototype.toString.call(value));
};

export const todayStr = () => {
	const today = new Date();
	return (
		today.getFullYear() +
		'-' +
		('0' + (today.getMonth() + 1)).slice(-2) +
		'-' +
		('0' + today.getDate()).slice(-2)
	);
};

export const weekAgoStr = () => {
	const today = new Date();
	const ago = new Date(today.getFullYear(), today.getMonth(), today.getDate() - 7);
	return (
		ago.getFullYear() +
		'-' +
		('0' + (ago.getMonth() + 1)).slice(-2) +
		'-' +
		('0' + ago.getDate()).slice(-2)
	);
};

export const monthAgoStr = () => {
	const today = new Date();
	const ago = new Date(today.getFullYear(), today.getMonth() - 1, today.getDate());
	return (
		ago.getFullYear() +
		'-' +
		('0' + (ago.getMonth() + 1)).slice(-2) +
		'-' +
		('0' + ago.getDate()).slice(-2)
	);
};
