import { Layout } from 'antd';
import {GithubOutlined}  from '@ant-design/icons'; 
const { Footer} = Layout;


export function PageFooter (props) {
    return(
        <Footer className={props.className} style={props.style}>
            ©2020 Created by <a href="https://robrotheram.com">Robrotheram</a> {window.w2g_version}  | Code Publicly avalible <a href="https://github.com/robrotheram/go-watch2gether"><GithubOutlined /></a>
        </Footer>
    )
}