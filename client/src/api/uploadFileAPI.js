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

const UploadFileInChunks = async (file, fileHash, chunksHash, uploadedChunksHash, token, addProgress) => {
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

		const sendRequest = async (chunkNum) => {
			const start = chunkNum * chunkSize;
			const end = ((start + chunkSize) >= file.size) ? file.size : start + chunkSize;
			const chunk = file.slice(start, end);

			const formData = new FormData();
			formData.append('chunk', chunk);
			formData.append('chunk_info', JSON.stringify({
				'file_name': file.name,
				'file_size': file.size,
				'file_hash': fileHash,
				'total_chunks': totalChunks,
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
				if (resp.data.file_success) {
					resolve(resp);
				}
				if (resp.data.chunk_success) {
					addProgress(Math.ceil((end - start) / file.size * 100));
					if (chunkNum + windowSize < totalChunks) {
						sendRequest(chunkNum + windowSize);
					}
				}
			} catch (err) {
				reject(err);
			}
		};

		for (let i = 0; i < Math.min(totalChunks, windowSize); i++) {
			sendRequest(i);
		}
	});
};

const UploadFileAPI = {
	CheckFileExistence,
	UploadFileInChunks,
};

export default UploadFileAPI;