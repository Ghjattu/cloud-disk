import React from 'react';
import './AuthenticationForm.css';
import { useState } from 'react';
import LoginForm from '../LoginForm/LoginForm';
import RegisterForm from '../RegisterForm/RegisterForm';

const AuthenticationForm = () => {
	const [isLoginForm, setIsLoginForm] = useState(true);

	const toggleForm = () => {
		setIsLoginForm(!isLoginForm);
	};

	return (
		<div className='authentication-form-container'>
			<div className='authentication-form'>
				{isLoginForm ? <LoginForm /> : <RegisterForm />}
			</div>
			<button onClick={toggleForm}>
				{isLoginForm ? 'Switch to Register' : 'Switch to Login'}
			</button>
		</div>
	);
};

export default AuthenticationForm;