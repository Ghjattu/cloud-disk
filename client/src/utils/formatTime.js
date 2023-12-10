const FormatTime = (time) => {
	const date = new Date(time * 1000);

	const now = new Date();
	const yesterday = new Date(now);
	yesterday.setDate(now.getDate() - 1);

	const isToday = date.toDateString() === now.toDateString();
	const isYesterday = date.toDateString() === yesterday.toDateString();

	if (isToday) {
		return `Today ${format(date)}`;
	} else if (isYesterday) {
		return `Yesterday ${format(date)}`;
	} else {
		// yy-mm-dd
		const year = date.getFullYear().toString().slice(-2);
		const month = date.getMonth() + 1;
		const day = date.getDate();
		return `${year}-${month}-${day}`;
	}
};

// return hh:mm
const format = (date) => {
	const hours = date.getHours();
	const minutes = padZero(date.getMinutes());

	return `${hours}:${minutes}`;
};

const padZero = (num) => {
	return num < 10 ? `0${num}` : num;
};

export default FormatTime;