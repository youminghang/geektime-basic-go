import dynamic  from 'next/dynamic'
import {Button, Form, Input} from "antd";
import {useEffect, useState} from "react";
import moment from "moment/moment";
import axios from "@/axios/axios";
import router from "next/router";
const WangEditor = dynamic(
    // 引入对应的组件 设置的组件参考上面的wangEditor react使用文档
    () => import('../../components/editor'),
    {ssr: false},
)

function Page() {
    const [html, setHtml] = useState('hello')

    const onFinish = (values: any) => {
        values.content = html
        axios.post("/articles/edit", values)
            .then((res) => {
                if(res.status != 200) {
                    alert(res.statusText);
                    return
                }
                if (res.data?.code == 0) {
                    router.push('/articles/'+ res.data?.data)
                    return
                }
                alert(res.data?.msg || "系统错误");
            }).catch((err) => {
            alert(err);
        })
    };

    const [data, setData] = useState<Article>( {id: 0, title: "", content: ""})
    const [isLoading, setLoading] = useState(false)
    useEffect(() => {
        const artID = router.query["id"]
        if (artID == undefined) {
            return
        }
        setLoading(true)
        axios.get('/articles/'+artID)
            .then((res) => res.data)
            .then((data) => {
                setData(data)
                setHtml(data.content)
                setLoading(false)
            })
    }, [])

    return <>
        <Form onFinish={onFinish} initialValues={{
            title: data.title
        }}>
            <Form.Item>
                <Input name={"title"} placeholder={"请输入标题"}/>
            </Form.Item>
            <WangEditor html={html} setHtmlFn={setHtml}/>
            <Form.Item>
                <Button type={"primary"}>保存</Button>
            </Form.Item>
        </Form>
    </>
}
export default Page