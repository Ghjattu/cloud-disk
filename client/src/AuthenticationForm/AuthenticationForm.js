import React, { useEffect } from 'react';
import './AuthenticationForm.css';
import { useState } from 'react';
import PropTypes from 'prop-types';
import LoginForm from '../LoginForm/LoginForm.js';
import RegisterForm from '../RegisterForm/RegisterForm.js';
import OAuthAPI from '../api/oauth.js';

const AuthenticationForm = ({ handleTokenChange }) => {
	const [isLoginForm, setIsLoginForm] = useState(true);

	useEffect(() => {
		const urlParams = new URLSearchParams(window.location.search);
		const code = urlParams.get('code');
		console.log(code);

		if (code) {
			handleGithubCallback(code);
		}
	}, []);

	const handleGithubCallback = async (code) => {
		try {
			const resp = await OAuthAPI.LoginWithGithub(code);
			if (resp.code !== 0) {
				throw new Error(resp.msg);
			}

			handleTokenChange(resp.data.token);
		} catch (err) {
			alert(err);
		}
	};

	const handleLoginWithGithub = async (e) => {
		e.preventDefault();

		try {
			const resp = await OAuthAPI.GetGithubLoginURL();
			if (resp.code !== 0) {
				throw new Error(resp.msg);
			}

			window.location.href = resp.data.url;
		} catch (err) {
			alert(err);
		}
	};

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
			<button onClick={handleLoginWithGithub}>Login with Github</button>
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