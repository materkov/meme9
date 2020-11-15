import React from 'react';
import * as schema from '../../schema/api';
import {Header} from "../Header/Header";
import {Link} from "../Link/Link";

export interface LoginPageProps {
    data: schema.LoginPageRenderer;
}

export class LoginPage extends React.Component<LoginPageProps> {
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

                {!data.headerRenderer?.currentUserId &&
                <Link onClick={this.onVKClick} href={data.vkUrl}>{data.vkText}</Link>
                }

                {data.headerRenderer?.currentUserId &&
                <span>
                    Вы уже вошли в систему
                    как <b>{data.headerRenderer.currentUserName}</b>. <Link href={"/logout"} onClick={this.logout}>Выйти</Link>
                </span>
                }

            </div>
        );
    }

}
