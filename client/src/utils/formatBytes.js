const FormatBytes = (bytes, decimals = 1) => {
	const units = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
	let unitIndex = 0;

	while (bytes >= 1024){
		bytes /= 1024;
		unitIndex++;
	}

	return `${bytes.toFixed(decimals)} ${units[unitIndex]}`;
};

export default FormatBytes;