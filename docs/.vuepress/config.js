module.exports = {
  title: "collection库", // 设置网站标题
  description: "一个让业务开发效率提高10倍的golang库", //描述
  dest: "../../dist/", // 设置输出目录
  port: 2333, //端口
  head: [["link", { rel: "icon", href: "/assets/img/head.png" }]],
  themeConfig: {
    // 添加导航栏
    nav: [
      { text: "主页", link: "/" }, // 导航条
      { text: "指南", link: "/guide/install" },
      { text: "手册", link: "/method/append" },
      {
        text: "github",
        // 这里是下拉列表展现形式。
        items: [
          {
            text: "collection",
            link: "https://github.com/goxstdx/go-collection",
          },
        ],
      },
    ],
    // 为以下路由添加侧边栏
    sidebar: {
      "/method/": [
        {
          title: "手册",
          collapsable: false,
          children: [
            "append",
            "avg",
            "contains",
            "filter",
            "map",
            "sum"
          ],
        },
      ],
      "/guide/": [
        {
          title: "指南",
          collapsable: true,
          children: [
            "install",
            "api",
            "usage-slice",
            "usage-map"
          ],
        },
      ],
    },
  },
};
