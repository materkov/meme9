import React from 'react';
import * as schema from '../../schema/login';
import {Header} from "../Header/Header";

interface State {
    login: string;
    password: string;
}

export class LoginPage extends React.Component<schema.LoginPageRenderer, State> {
    state: State = {
        login: '',
        password: '',
    };

    onLoginChanged = (e: React.ChangeEvent<HTMLInputElement>) => {
        this.setState({login: e.target.value});
    };

    onPasswordChanged = (e: React.ChangeEvent<HTMLInputElement>) => {
        this.setState({password: e.target.value});
    };

    onSubmit = () => {
        const params: schema.LoginRequest = {
            login: this.state.login,
            password: this.state.password,
        };

        fetch("/api/login", {
            method: 'POST', body: JSON.stringify(params)
        })
    };

    render() {
        return (
            <div>
                {this.props.headerRenderer && <Header {...this.props.headerRenderer}/>}

                {this.props.welcomeText}
                <br/>

                <input type="text" placeholder="Логин" value={this.state.login} onChange={this.onLoginChanged}/>
                <br/>
                <input type="password" placeholder="Пароль" value={this.state.password}
                       onChange={this.onPasswordChanged}/>
                <br/>
                <button onClick={this.onSubmit}>Войти</button>
            </div>
        );
    }
}
