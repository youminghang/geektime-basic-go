import dynamic  from 'next/dynamic'
import {Button, Form, Input} from "antd";
import {useEffect, useState} from "react";
import axios from "@/axios/axios";
import router from "next/router";
import {ProLayout} from "@ant-design/pro-components";
import {useSearchParams} from "next/navigation";
import {useForm} from "antd/es/form/Form";
const WangEditor = dynamic(
    // 引入对应的组件 设置的组件参考上面的wangEditor react使用文档
    () => import('../../components/editor'),
    {ssr: false},
)

function Page() {
    const [html, setHtml] = useState()
    const params = useSearchParams()
    const artID = params?.get("id")
    const onFinish = (values: any) => {
        if(artID) {
            values.id = parseInt(artID)
        }
        values.content = html
        axios.post("/articles/edit", values)
            .then((res) => {
                if(res.status != 200) {
                    alert(res.statusText);
                    return
                }
                if (res.data?.code == 0) {
                    router.push('/articles/list')
                    return
                }
                alert(res.data?.msg || "系统错误");
            }).catch((err) => {
            alert(err);
        })
    };

    const [data, setData] = useState<Article>( {id: 0, title: "", content: ""})
    const [form] = Form.useForm()
    useEffect(() => {
        if (!artID) {
            return
        }
        axios.get('/articles/detail/'+artID)
            .then((res) => res.data)
            .then((data) => {
                form.setFieldsValue(data.data)
                // setData()
                console.log(data)
                setHtml(data.data.content)
            })
    }, [artID])

    return <ProLayout title={"创作中心"}>
        <Form onFinish={onFinish}
        form={form}
              initialValues={data}>
            <Form.Item name={"title"}
                       rules={[{ required: true, message: '请输入标题' }]}
            >
                <Input placeholder={"请输入标题"}/>
            </Form.Item>
            <WangEditor html={html} setHtmlFn={setHtml}/>
            <Form.Item>
                <br/>
                <Button type={"primary"} htmlType={"submit"}>保存</Button>
            </Form.Item>
        </Form>
    </ProLayout>
}
export default Page