import { Layout } from 'antd';
import {GithubOutlined}  from '@ant-design/icons'; 
const { Footer} = Layout;


export function PageFooter () {
    return(
        <Footer style={{ textAlign: 'center', position: "fixed", bottom:"0px", left:"0px", right:"0px", height:"50px", padding: "15px 50px 28px 50px" }}>
            ©2020 Created by <a href="https://robrotheram.com">Robrotheram</a> {window.w2g_version}  | Code Publicly avalible <a href="https://github.com/robrotheram/go-watch2gether"><GithubOutlined /></a>
        </Footer>
    )
}