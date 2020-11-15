import React from 'react';
import * as schema from '../../schema/api';
import {Header} from "../Header/Header";
import {Link} from "../Link/Link";

export interface LoginPageProps {
    data: schema.LoginPageRenderer;
}

interface LoginPageState {
    login: string;
    password: string;
    response?: schema.LoginRenderer;
    logoutResponse?: schema.LogoutRenderer;
}

export class LoginPage extends React.Component<LoginPageProps, LoginPageState> {
    state: LoginPageState = {
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
        this.setState({response: undefined});

        const params: schema.LoginRequest = {
            login: this.state.login,
            password: this.state.password,
        };

        fetch("/api/login", {
            method: 'POST', body: JSON.stringify(params)
        }).then(r => r.json()).then((r: schema.LoginRenderer) => {
            this.setState({response: r});
        })
    };

    logout = () => {
        window.location.href = "/logout";
    }

    onVKClick = () => {
        window.location.href = this.props.data.vkUrl;
    };

    render() {
        const data = this.props.data;

        return (
            <div>
                <Header data={data.headerRenderer}/>

                <h2>{data.welcomeText}</h2>
                <br/>

                Войти через <Link onClick={this.onVKClick} href={""}>VK</Link><br/>

                {data.headerRenderer?.currentUserId ?
                    <span>
                        Вы уже вошли в систему
                        как <b>{data.headerRenderer.currentUserName}</b>. <Link href={""}
                                                                                      onClick={this.logout}>Выйти</Link>
                    </span> :
                    this.renderForm()
                }

                {this.state.logoutResponse &&
                <LogoutSuccess/>
                }
            </div>
        );
    }

    renderForm() {
        return (
            <>
                <input type="text" placeholder="Логин" value={this.state.login} onChange={this.onLoginChanged}/>
                <br/>
                <input type="password" placeholder="Пароль" value={this.state.password}
                       onChange={this.onPasswordChanged}/>
                <br/><br/>
                <button onClick={this.onSubmit}>Войти</button>
            </>
        )
    }
}

interface ErrorProps {
    data: schema.ErrorRenderer;
}

function Error(props: ErrorProps) {
    return <div style={{"background": "red", "padding": "10px"}}>
        {props.data.displayText}
    </div>;
}

interface SuccessProps {
    data: schema.LoginRenderer
}

function Success(props: SuccessProps) {
    return <div>Вы вошли как {props.data.headerRenderer?.currentUserName}. Обновите страницу.</div>
}

interface LogoutSuccessProps {
}

function LogoutSuccess() {
    return <div>Вы вышли из системы. Обновите странцу.</div>
}
