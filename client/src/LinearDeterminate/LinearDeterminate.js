import React from 'react';
import Box from '@mui/material/Box';
import LinearProgress from '@mui/material/LinearProgress';
import PropTypes from 'prop-types';

const LinearDeterminate = ({ progress }) => {
	return (
		<Box sx={{ width: '100%' }}>
			<LinearProgress
				variant="determinate"
				value={progress}
				style={{ height: 7, borderRadius: 5 }}
			/>
		</Box>
	);
};

LinearDeterminate.propTypes = {
	progress: PropTypes.number.isRequired,
};

export default LinearDeterminate;