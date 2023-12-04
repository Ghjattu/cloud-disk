import React from 'react';
import './AuthenticationForm.css';
import { useState } from 'react';
import PropTypes from 'prop-types';
import LoginForm from '../LoginForm/LoginForm';
import RegisterForm from '../RegisterForm/RegisterForm';

const AuthenticationForm = ({ handleTokenChange }) => {
	const [isLoginForm, setIsLoginForm] = useState(true);

	const toggleForm = () => {
		setIsLoginForm(!isLoginForm);
	};

	return (
		<div className='authentication-form-container'>
			<div className='authentication-form'>
				{isLoginForm ?
					<LoginForm handleTokenChange={handleTokenChange} /> :
					<RegisterForm handleTokenChange={handleTokenChange} />}
			</div>
			<button onClick={toggleForm}>
				{isLoginForm ? 'Switch to Register' : 'Switch to Login'}
			</button>
		</div>
	);
};

AuthenticationForm.propTypes = {
	handleTokenChange: PropTypes.func.isRequired,
};

export default AuthenticationForm;