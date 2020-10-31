import React from 'react';
import * as schema from '../../schema/login';
import {AnyRenderer} from '../../schema/login';
import {Header} from "../Header/Header";
import {Link} from "../Link/Link";

interface LoginPageState {
    login: string;
    password: string;
    response?: schema.AnyRenderer;
    logoutResponse?: schema.AnyRenderer;
}

export class LoginPage extends React.Component<schema.LoginPageRenderer, LoginPageState> {
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
        }).then(r => r.json()).then((r: AnyRenderer) => {
            this.setState({response: r});
        })
    };

    logout = () => {
        const params: schema.LogoutRequest = {};
        fetch("/api/logout", {
            method: 'POST', body: JSON.stringify(params),
        }).then(r => r.json()).then((r: AnyRenderer) => {
            this.setState({logoutResponse: r});
        });
    }

    render() {
        return (
            <div>
                {this.props.headerRenderer && <Header data={this.props.headerRenderer}/>}

                <h2>{this.props.welcomeText}</h2>
                <br/>

                {this.props.headerRenderer?.currentUserId ?
                    <span>
                        Вы уже вошли в систему
                        как <b>{this.props.headerRenderer.currentUserName}</b>. <Link href={""}
                                                                                      onClick={this.logout}>Выйти</Link>
                    </span> :
                    this.renderForm()
                }

                {this.state.response?.errorRenderer &&
                <Error data={this.state.response.errorRenderer}/>
                }

                {this.state.response?.loginRenderer &&
                <Success data={this.state.response.loginRenderer}/>
                }

                {this.state.logoutResponse?.logoutRenderer &&
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
