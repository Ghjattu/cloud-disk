import SparkMD5 from 'spark-md5';

const CalculateMD5 = async (file) => {
	return new Promise((resolve, reject) => {
		const fileReader = new FileReader();

		fileReader.onload = () => {
			const result = fileReader.result;
			const fileHash = SparkMD5.ArrayBuffer.hash(result);

			resolve(fileHash);
		};

		fileReader.onerror = () => {
			reject(new Error('File reading error'));
		};

		fileReader.readAsArrayBuffer(file);
	});
};

export default CalculateMD5;