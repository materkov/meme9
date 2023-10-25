import fetch from "node-fetch";
import {XMLParser} from "fast-xml-parser";
import mysql from "mysql2";

let token = '';

function fetchRss() {
    return new Promise((resolve, reject) => {
        fetch('https://rss.nytimes.com/services/xml/rss/nyt/World.xml')
            .then(r => {
                if (!r.ok) {
                    reject();
                    return;
                }

                r.text()
                    .then(text => resolve(text))
                    .catch(reject)
            })
            .catch(reject);
    });
}

function parse() {
    console.log(new Date().toUTCString() + ' Starting parsing');

    apiCall('posts.listPostedByUser', {userId: "77"}).then(result => {
        console.log(new Date().toUTCString() + ' Fetched users posts');

        let parsedUrls = [];
        const urlRegex = /(https?:\/\/[^\s]+)/g;
        for (let post of result) {
            const matches = post.text.match(urlRegex);
            if (matches) {
                parsedUrls = [...parsedUrls, ...matches];
            }
        }

        fetchRss()
            .then(text => {
                console.log(new Date().toUTCString() + ' Fetched RSS feed');

                const parser = new XMLParser();
                let rss = parser.parse(text);
                const items = rss['rss']['channel']['item'].slice(0, 5);
                for (let item of items) {
                    if (parsedUrls.indexOf(item['link']) !== -1) {
                        console.log(new Date().toUTCString() + ' Skipped item ', item['link']);
                        continue;
                    }

                    console.log(new Date().toUTCString() + ' Adding item ', item['link']);
                    const text = item['title'] + "." + "\n" + item['description'] + "\n\n" + item['link'];

                    apiCall('posts.add', {
                        text: text,
                    }, {
                        Authorization: 'Bearer ' + token,
                    })
                        .then(() => {
                        })
                        .catch(() =>
                            console.log()
                        );
                }
            })
    })


}

function apiCall(method, args, meta) {
    return new Promise((resolve, reject) => {
        fetch("https://meme.mmaks.me/api/" + method, {
            method: 'POST',
            body: JSON.stringify(args || {}),
            headers: meta || [],
        })
            .then(r => r.json())
            .then(r => {
                if (r.error) {
                    reject(r.error);
                } else {
                    resolve(r);
                }
            })
            .catch(reject);
    })
}

function getConfig() {
    return new Promise((resolve, reject) => {
        const conn = mysql.createConnection({
            host: "localhost",
            user: "root",
            password: "root",
            database: "meme9",
        });

        conn.connect(function (err) {
            if (err) {
                reject(err);
                return;
            }

            conn.query("SELECT * FROM objects where id = -5", function (err, result, fields) {
                if (err) {
                    reject(err);
                    return;
                }

                resolve(result[0].data['RSSToken']);
            });
        });
    });
}

getConfig().then(t => {
    token = t

    parse();
    setInterval(() => {
        parse();
    }, 60 * 1000);
})
