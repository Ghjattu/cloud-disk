const CheckFileExistence = async (fileHash, token) => {
	const url = `/file/exist/${fileHash}`;
	const headers = {
		'Content-Type': 'application/json',
		'Accept': 'application/json',
		'Authorization': `Bearer ${token}`,
	};

	const response = await fetch(url, {
		method: 'GET',
		headers: headers,
	});
	return await response.json();
};

const UploadChunks = async (file, fileHash, chunksHash, uploadedChunksHash, token, addProgress, controller) => {
	const url = '/file/upload';
	const headers = {
		// add this line will cause error
		// 'Content-Type': 'multipart/form-data',
		'Authorization': `Bearer ${token}`,
	};

	return new Promise((resolve, reject) => {
		// eslint-disable-next-line no-undef
		const windowSize = parseInt(process.env.REACT_APP_WINDOW_SIZE);
		// eslint-disable-next-line no-undef
		const chunkSize = parseInt(process.env.REACT_APP_CHUNK_SIZE);
		const totalChunks = Math.ceil(file.size / chunkSize);
		let successfulChunks = 0;

		const sendRequest = async (chunkNum) => {
			const start = chunkNum * chunkSize;
			const end = ((start + chunkSize) >= file.size) ? file.size : start + chunkSize;
			const chunk = file.slice(start, end);

			const formData = new FormData();
			formData.append('chunk', chunk);
			formData.append('chunk_info', JSON.stringify({
				'file_hash': fileHash,
				'chunk_hash': chunksHash[chunkNum],
				'chunk_num': chunkNum,
			}));

			try {
				const response = await fetch(url, {
					method: 'POST',
					headers: headers,
					body: formData,
				});
				const resp = await response.json();
				if (resp.code !== 0) {
					throw new Error(resp.msg);
				}
				if (resp.data.chunk_success) {
					successfulChunks++;
					addProgress(Math.ceil((end - start) / file.size * 100));
					if (successfulChunks === totalChunks) {
						resolve();
					} else if (chunkNum + windowSize < totalChunks) {
						sendRequest(chunkNum + windowSize);
					}
				}
			} catch (err) {
				controller.abort();
				reject(err);
			}
		};

		for (let i = 0; i < Math.min(totalChunks, windowSize); i++) {
			sendRequest(i);
		}
	});
};

const MergeChunks = async (fileHash, fileName, fileSize, token) => {
	const url = '/file/merge';
	const headers = {
		'Content-Type': 'application/json',
		'Authorization': `Bearer ${token}`,
	};
	const data = JSON.stringify({
		'file_hash': fileHash,
		'file_name': fileName,
		'file_size': fileSize,
	});

	return new Promise((resolve, reject) => {
		const sendRequest = async () => {
			try {
				const response = await fetch(url, {
					method: 'POST',
					headers: headers,
					body: data,
				});
				const resp = await response.json();
				if (resp.code !== 0) {
					throw new Error(resp.msg);
				}
				resolve(resp.data);
			} catch (err) {
				reject(err);
			}
		};

		sendRequest();
	});
};

const UploadFileAPI = {
	CheckFileExistence,
	UploadChunks,
	MergeChunks,
};

export default UploadFileAPI;