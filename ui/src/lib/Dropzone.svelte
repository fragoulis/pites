<script lang="ts">
	import { Dropzone } from 'flowbite-svelte';
	import pb from '$lib/pocketbase';
	import { sharedToast } from '$lib/store';

	let fileToUpload: File;
	let importing = false;
	let importError: any;

	const supportedFileTypes: string[] = [
		'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
	];

	const dropHandle = (event: any) => {
		event.preventDefault();

		if (event.dataTransfer.items) {
			[...event.dataTransfer.items].forEach((item, i) => {
				uploadFile(item.getAsFile());
			});
		} else {
			[...event.dataTransfer.files].forEach((file, i) => {
				uploadFile(file);
			});
		}
	};

	const handleChange = (event: any) => {
		const files = event.target.files;
		if (files.length > 0) {
			uploadFile(files[0]);
		}
	};

	const showFiles = () => {
		return fileToUpload.name;
	};

	const uploadFile = async (file: File) => {
		if (!supportedFileTypes.includes(file.type)) {
			console.debug('file type is no supported: ', file.type);

			$sharedToast.show = true;
			$sharedToast.success = false;
			$sharedToast.message = 'Ο τύπος αρχείου δεν υποστηρίζεται.';

			return;
		}

		fileToUpload = file;

		const formData = new FormData();
		formData.append('file', fileToUpload);

		try {
			const response = await pb.send('/members/import', {
				method: 'POST',
				body: formData
			});

			console.log(response);

			importing = false;
			importError = undefined;

			$sharedToast.show = true;
			$sharedToast.success = true;
			$sharedToast.message = 'Η εισαγωγή ολοκληρώθηκε με επιτυχία.';
		} catch (error: any) {
			importError = error;

			$sharedToast.show = true;
			$sharedToast.success = false;
			$sharedToast.message = 'Παρουσιάστηκε κάποιο πρόβλημα.';

			importing = false;
		}
	};
</script>

<Dropzone
	id="dropzone"
	on:drop={dropHandle}
	on:dragover={(event) => {
		event.preventDefault();
	}}
	on:change={handleChange}
>
	<svg
		aria-hidden="true"
		class="mb-3 w-10 h-10 text-gray-400"
		fill="none"
		stroke="currentColor"
		viewBox="0 0 24 24"
		xmlns="http://www.w3.org/2000/svg"
		><path
			stroke-linecap="round"
			stroke-linejoin="round"
			stroke-width="2"
			d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
		/></svg
	>

	{#if !fileToUpload}
		<p class="mb-2 text-sm text-gray-500 dark:text-gray-400">
			<span class="font-semibold">Κάντε κλικ για να ανεβάσετε</span>
			ή τραβήξτε το αρχείο σε αυτό το χώρο.
		</p>
		<p class="text-xs text-gray-500 dark:text-gray-400">.xlsx</p>
	{:else if importing}
		<p>Η εισαγωγή ξεκίνησε...</p>
	{:else if importError}
		<p>Η εισαγωγή απέτυχε.</p>
		<p>{importError}</p>
	{:else}
		<p>Η εισαγωγή ολοκληρώθηκε.</p>
	{/if}
</Dropzone>
