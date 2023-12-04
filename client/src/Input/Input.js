import React from 'react';
import './Input.css';
import PropTypes from 'prop-types';

const Input = ({ label, name, value, onChange }) => {
	return (
		<div>
			<label htmlFor={name}>{label}</label>
			<input type='text' id={name} name={name} value={value} onChange={onChange} />
		</div>
	);
};

Input.propTypes = {
	label: PropTypes.string.isRequired,
	name: PropTypes.string.isRequired,
	value: PropTypes.string.isRequired,
	onChange: PropTypes.func.isRequired,
};

export default Input;