import SparkMD5 from 'spark-md5';
import GetChunkSize from './getChunkSize.js';

// Calculate MD5 hash of file in chunks
const CalculateMD5 = async (file) => {
	return new Promise((resolve, reject) => {
		const chunkSize = GetChunkSize();
		const totalChunks = Math.ceil(file.size / chunkSize);
		let currentChunk = 0;
		const sparkArrayBuffer = new SparkMD5.ArrayBuffer();
		const fileReader = new FileReader();

		fileReader.onload = () => {
			console.log('currentChunk: ', currentChunk);
			sparkArrayBuffer.append(fileReader.result);
			currentChunk++;

			if (currentChunk < totalChunks) {
				loadNext();
			} else {
				console.log('finished loading file');
				const fileHash = sparkArrayBuffer.end();
				console.log('fileHash: ', fileHash);

				resolve(fileHash);
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