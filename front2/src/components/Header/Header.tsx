import React from "react";
import styles from "./Header.module.css";
import {Link} from "../Link/Link";
import {FeedGetHeaderRequest, FeedGetHeaderResponse, HeaderRenderer} from "../../api/api2";
import {api} from "../../Api";

interface State {
    data: HeaderRenderer;
}

export class Header extends React.Component {
    state: State = {
        data: HeaderRenderer.fromPartial({}),
    };

    componentDidMount() {
        this.refreshData();
        setInterval(this.refreshData, 60 * 1000);
    }

    refreshData = () => {
        api<FeedGetHeaderRequest, FeedGetHeaderResponse>("meme.Feed.GetHeader", {}).then(r => {
            this.setState({data: r.renderer});
        }).catch(() => {
            console.error('Failed updating header');
        })
    }

    render() {
        const data = this.state.data;

        return (
            <div className={styles.Header}>
                <Link href={data.mainUrl} className={styles.Logo}>meme</Link>

                <div className={styles.RightContainer}>
                    {data.isAuthorized && <div className={styles.Name}>{data.userName}</div>}
                    {data.isAuthorized && <img className={styles.Avatar} alt="" src={data.userAvatar}/>}
                    {data.isAuthorized &&
                    <a href={data.logoutUrl}>Выход</a>
                    }

                    {!data.isAuthorized && <Link className={styles.Name} href={"/login"}>Войти</Link>}
                </div>
            </div>
        );
    }
}
