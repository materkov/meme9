import React from "react";
import * as schema from "../../schema/login";
import {Header} from "../Header/Header";

interface State {
    text: string;
}

export class Composer extends React.Component<schema.ComposerRenderer, State> {
    state: State = {
        text: '',
    };

    onChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        this.setState({text: e.target.value});
    };

    onSubmit = () => {
        //@ts-ignore
        const params: schema.AnyRequest = {
            addPostRequest: {
                text: this.state.text,
            }
        };

        fetch("/api", {
            method: 'POST',
            body: JSON.stringify(params),
        }).then(r => r.json()).then(r => {

        });
    };

    render() {
        return (
            <div>
                {this.props.headerRenderer && <Header {...this.props.headerRenderer}/>}

                Напишите свой пост здесь:<br/>
                {this.props.welcomeText}<br/>
                <textarea onChange={this.onChange}></textarea><br/>
                <button onClick={this.onSubmit}>Отправить</button>
            </div>
        );
    }
}
