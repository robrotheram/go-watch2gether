import { Layout } from 'antd';
const { Footer} = Layout;

export function PageFooter () {
    return(
        <Footer style={{ textAlign: 'center', position: "fixed", bottom:"0px", left:"0px", right:"0px", height:"50px", padding: "15px 50px 28px 50px" }}>
            ©2020 Created by Robrotheram {window.w2g_version}
        </Footer>
    )
}