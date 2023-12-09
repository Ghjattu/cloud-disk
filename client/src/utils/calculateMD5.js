import SparkMD5 from 'spark-md5';

// Calculate MD5 hash of file in chunks
const CalculateMD5 = async (file) => {
	return new Promise((resolve, reject) => {
		// eslint-disable-next-line no-undef
		const chunkSize = parseInt(process.env.REACT_APP_CHUNK_SIZE);
		const totalChunks = Math.ceil(file.size / chunkSize);
		const chunksHash = [];
		let currentChunk = 0;
		const sparkArrayBuffer = new SparkMD5.ArrayBuffer();
		const fileReader = new FileReader();

		fileReader.onload = () => {
			sparkArrayBuffer.append(fileReader.result);
			chunksHash.push(SparkMD5.ArrayBuffer.hash(fileReader.result));
			currentChunk++;

			if (currentChunk < totalChunks) {
				loadNext();
			} else {
				const fileHash = sparkArrayBuffer.end();

				resolve({
					fileHash,
					chunksHash,
				});
			}

		};

		fileReader.onerror = () => {
			reject(new Error('File reading error'));
		};

		const loadNext = () => {
			const start = currentChunk * chunkSize;
			const end = ((start + chunkSize) >= file.size) ? file.size : start + chunkSize;

			fileReader.readAsArrayBuffer(file.slice(start, end));
		};

		loadNext();
	});
};

export default CalculateMD5;