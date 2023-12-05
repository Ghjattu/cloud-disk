import GetChunkSize from '../utils/getChunkSize.js';

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

const UploadFileInChunks = async (file, fileHash, chunksHash, uploadedChunksHash, token) => {
	console.log('start uploading file in chunks');
	const url = '/file/upload';
	const headers = {
		// add this line will cause error
		// 'Content-Type': 'multipart/form-data',
		'Authorization': `Bearer ${token}`,
	};

	return new Promise((resolve, reject) => {
		const chunkSize = GetChunkSize(); // 100 KB
		const totalChunks = Math.ceil(file.size / chunkSize);

		for (let chunkNum = 0; chunkNum < totalChunks; chunkNum++) {
			if (uploadedChunksHash.includes(chunksHash[chunkNum])) {
				continue;
			}

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
				'chunk_num': chunkNum
			}));

			fetch(url, {
				method: 'POST',
				headers: headers,
				body: formData,
			})
				.then((response) => response.json())
				.then((resp) => {
					console.log('Upload chunk %d successfully', chunkNum);
					if (resp.file_success) {
						resolve(resp);
					}
				})
				.catch((error) => {
					reject(error);
				});
		}
	});
};

const UploadFileAPI = {
	CheckFileExistence,
	UploadFileInChunks,
};

export default UploadFileAPI;