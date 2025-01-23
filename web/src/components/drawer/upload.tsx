import { useState } from "react";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { useForm, useFieldArray } from "react-hook-form";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { SheetClose } from "@/components/ui/sheet";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import * as http from "@/api/modules/fouceinter";
import { useTranslation } from "react-i18next";
import { ChevronLeftIcon, PlusIcon, TrashIcon, FilePlusIcon } from "@radix-ui/react-icons";
import { 
  Select, 
  SelectContent, 
  SelectItem, 
  SelectTrigger, 
  SelectValue 
} from "@/components/ui/select";

interface UploadFormValues {
  class: string;
  command: string;
  depends_version: string;
  description: string;
  docker_compose: string;
  env: Env[];
  github: string;
  icon: string;
  key: string;
  name: string;
  nginx_config: string;
  repo: string;
  version: string;
  volume: Volume[];
}

interface Env {
  key: string;
  name: string;
  type: string;
  value: string;
  required: boolean;
}

interface Volume {
  local: string;
  target: string;
}

interface UploadSheetProps {
  isOpen: boolean;
  onClose: () => void;
}

export function UploadSheet({ isOpen, onClose }: UploadSheetProps) {
  const { t } = useTranslation();
  const [loading, setLoading] = useState<boolean>(false);

  const defaultValues: UploadFormValues = {
    class: "",
    command: "",
    depends_version: "",
    description: "",
    docker_compose: "",
    env: [],
    github: "",
    icon: "",
    key: "",
    name: "",
    nginx_config: "",
    repo: "",
    version: "",
    volume: [],
  };

  const form = useForm<UploadFormValues>({
    defaultValues,
  });

  const {
    control,
    handleSubmit,
    register,
    setValue,
    formState: { errors },
  } = form;

  const { fields: envFields, append: appendEnv, remove: removeEnv } = useFieldArray({
    control,
    name: "env",
  });

  const { fields: volumeFields, append: appendVolume, remove: removeVolume } = useFieldArray({
    control,
    name: "volume",
  });

  const handleUpload = async (data: UploadFormValues) => {
    setLoading(true);
    
    try {
      const response = await http.uploadApp(data);
      if (response.code === 200) {
        onClose();
      }
      console.log("上传应用", data)
    } catch (error) {
      console.error("Upload failed:", error);
    } finally {
      setLoading(false);
    }
  };

   // 生成默认 Docker Compose 配置
   const generateDefaultDockerCompose = () => {
    const defaultDockerCompose = `services:
  ${form.getValues('key')}:
    image: ${form.getValues('repo') || 'unknown'}:${form.getValues('version') || 'latest'}
    restart: always
    container_name: \${CONTAINER_NAME}
    networks:
      \${DOOTASK_NETWORK_NAME}:
        ipv4_address: \${IP_ADDRESS}
    environment:
      - CLOUD_PROVIDER=\${CLOUD_PROVIDER}
    cpus: "\${CPUS}"
    mem_limit: "\${MEMORY_LIMIT}"
    labels:
      createdBy: "Apps"

networks:
  \${DOOTASK_NETWORK_NAME}:
    external: true

`;
    setValue('docker_compose', defaultDockerCompose);
  };

  // 生成默认 Nginx 配置
  const generateDefaultNginxConfig = () => {
    const defaultNginxConfig = `location /plugin/{{.Key}}/ {
    proxy_http_version 1.1;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Real-PORT $remote_port;
    proxy_set_header X-Forwarded-Host $the_host;
    proxy_set_header X-Forwarded-Proto $the_scheme;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header Host $http_host;
    proxy_set_header Scheme $scheme;
    proxy_set_header Server-Protocol $server_protocol;
    proxy_set_header Server-Name $server_name;
    proxy_set_header Server-Addr $server_addr;
    proxy_set_header Server-Port $server_port;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection $connection_upgrade;
    proxy_pass http://{{.ContainerName}}:{{.Port}}/;
}`;
    setValue('nginx_config', defaultNginxConfig);
  };

  const renderDynamicEnvInputs = () => (
    <div>
      <div className="flex items-center justify-between mb-2">
        <FormLabel>{t("环境变量")}</FormLabel>
        <Button 
          type="button" 
          variant="surely"
          size="sm" 
          onClick={() => appendEnv({ key: "", name: "", type: "string", value: "", required: false })}
        >
          <PlusIcon className="mr-2" /> {t("添加环境变量")}
        </Button>
      </div>
      {envFields.map((field, index) => (
        <div key={field.id} className="grid grid-cols-6 gap-2 mb-2">
          <FormItem className="col-span-2">
            <FormControl>
              <Input
                {...register(`env.${index}.key` as const)}
                placeholder={t("变量键名")}
              />
            </FormControl>
          </FormItem>
          <FormItem className="col-span-2">
            <FormControl>
              <Input
                {...register(`env.${index}.name` as const)}
                placeholder={t("变量名称")}
              />
            </FormControl>
          </FormItem>
          <FormItem className="col-span-1">
            <Select
              {...register(`env.${index}.type` as const)}
              defaultValue="string"
            >
              <SelectTrigger>
                <SelectValue placeholder={t("类型")} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="string">String</SelectItem>
                <SelectItem value="number">Number</SelectItem>
                <SelectItem value="boolean">Boolean</SelectItem>
              </SelectContent>
            </Select>
          </FormItem>
          <Button 
            type="button" 
            variant="destructive" 
            size="icon" 
            onClick={() => removeEnv(index)}
          >
            <TrashIcon />
          </Button>
        </div>
      ))}
    </div>
  );

  const renderDynamicVolumeInputs = () => (
    <div>
      <div className="flex items-center justify-between mb-2">
        <FormLabel>{t("数据卷映射")}</FormLabel>
        <Button 
          type="button" 
          variant="surely"
          size="sm" 
          onClick={() => appendVolume({ local: "", target: "" })}
        >
          <PlusIcon className="mr-2" /> {t("添加数据卷")}</Button>
      </div>
      {volumeFields.map((field, index) => (
        <div key={field.id} className="grid grid-cols-5 gap-2 mb-2">
          <FormItem className="col-span-2">
            <FormControl>
              <Input
                {...register(`volume.${index}.local` as const)}
                placeholder={t("本地路径")}
              />
            </FormControl>
          </FormItem>
          <FormItem className="col-span-2">
            <FormControl>
              <Input
                {...register(`volume.${index}.target` as const)}
                placeholder={t("容器路径")}
              />
            </FormControl>
          </FormItem>
          <Button 
            type="button" 
            variant="destructive" 
            size="icon" 
            onClick={() => removeVolume(index)}
          >
            <TrashIcon />
          </Button>
        </div>
      ))}
    </div>
  );

  return (
    <Sheet open={isOpen} onOpenChange={onClose}>
      <SheetContent className="overflow-y-auto">
        <SheetHeader>
          <SheetTitle className="lg:ml-2 md:ml-2 text-gray-700 z-50 lg:bg-transparent md:bg-transparent bg-gray-200/50 lg:pb-1 md:py-0 py-3 flex items-center gap-2">
            <ChevronLeftIcon
              className="h-6 w-6 ml-1 lg:hidden md:hidden block"
              onClick={() => onClose()}
            />
            {t("应用上传")}
          </SheetTitle>
        </SheetHeader>
        <hr  className='lg:block md:block hidden'/>
        <Form {...form}>
          <form
            className="space-y-8 relative overflow-visible lg:px-0 md:px-0 px-3 pb-3 py-5"
            onSubmit={handleSubmit(handleUpload)}
          >
            {/* 基础表单项 */}
            <FormItem>
              <FormLabel>{t("应用名称")}</FormLabel>
              <FormControl>
                <Input
                  {...register("name", { required: t("应用名称不能为空") })}
                  className="bg-gray-200/60 border border-gray-200/60"
                  placeholder={t("请输入应用名称...")}
                />
              </FormControl>
              <FormMessage>{errors.name?.message}</FormMessage>
            </FormItem>

            <FormItem>
              <FormLabel>{t("应用标识")}</FormLabel>
              <FormControl>
                <Input
                  {...register("key")}
                  className="bg-gray-200/60 border border-gray-200/60"
                  placeholder={t("请输入应用唯一标识...")}
                />
              </FormControl>
            </FormItem>

            <FormItem>
              <FormLabel>{t("应用分类")}</FormLabel>
              <FormControl>
                <Input
                  {...register("class")}
                  className="bg-gray-200/60 border border-gray-200/60"
                  placeholder={t("请输入应用分类...")}
                />
              </FormControl>
            </FormItem>

            <FormItem>
              <FormLabel>{t("应用描述")}</FormLabel>
              <FormControl>
                <Textarea
                  {...register("description")}
                  className="bg-gray-200/60 border border-gray-200/60"
                  placeholder={t("请输入应用描述...")}
                />
              </FormControl>
            </FormItem>

            <div className="grid grid-cols-2 gap-4">
              <FormItem>
                <FormLabel>{t("镜像仓库")}</FormLabel>
                <FormControl>
                  <Input
                    {...register("repo")}
                    className="bg-gray-200/60 border border-gray-200/60"
                    placeholder={t("请输入镜像仓库...")}
                  />
                </FormControl>
              </FormItem>

              <FormItem>
                <FormLabel>{t("镜像版本")}</FormLabel>
                <FormControl>
                  <Input
                    {...register("version")}
                    className="bg-gray-200/60 border border-gray-200/60"
                    placeholder={t("请输入镜像版本...")}
                  />
                </FormControl>
              </FormItem>
            </div>

            <FormItem>
              <FormLabel>{t("依赖版本")}</FormLabel>
              <FormControl>
                <Input
                  {...register("depends_version")}
                  className="bg-gray-200/60 border border-gray-200/60"
                  placeholder={t("请输入依赖版本...")}
                />
              </FormControl>
            </FormItem>

            <FormItem>
              <FormLabel>{t("Github地址")}</FormLabel>
              <FormControl>
                <Input
                  {...register("github")}
                  className="bg-gray-200/60 border border-gray-200/60"
                  placeholder={t("请输入Github地址...")}
                />
              </FormControl>
            </FormItem>

            <FormItem>
              <FormLabel>{t("应用图标")}</FormLabel>
              <FormControl>
                <Input
                  {...register("icon")}
                  className="bg-gray-200/60 border border-gray-200/60"
                  placeholder={t("请输入图标地址或Base64...")}
                />
              </FormControl>
            </FormItem>

            <FormItem>
              <div className="flex justify-between items-center">
                <FormLabel>{t("Docker Compose配置")}</FormLabel>
                <Button 
                  type="button" 
                  variant="surely"
                  size="sm" 
                  onClick={generateDefaultDockerCompose}
                  disabled={!form.getValues('repo') || !form.getValues('key')}
                >
                  <FilePlusIcon className="mr-2" /> {t("生成默认配置")}
                </Button>
              </div>
              <FormControl>
                <Textarea
                  {...register("docker_compose")}
                  className="bg-gray-200/60 border border-gray-200/60"
                  placeholder={t("请输入Docker Compose配置...")}
                  rows={5}
                />
              </FormControl>
            </FormItem>

            <FormItem>
              <div className="flex justify-between items-center">
                <FormLabel>{t("Nginx配置")}</FormLabel>
                <Button 
                  type="button" 
                  variant="surely"
                  size="sm" 
                  onClick={generateDefaultNginxConfig}
                >
                  <FilePlusIcon className="mr-2" /> {t("生成默认配置")}
                </Button>
              </div>
              <FormControl>
                <Textarea
                  {...register("nginx_config")}
                  className="bg-gray-200/60 border border-gray-200/60"
                  placeholder={t("请输入Nginx配置...")}
                  rows={5}
                />
              </FormControl>
            </FormItem>

            {/* 动态环境变量表单项 */}
            {renderDynamicEnvInputs()}

            {/* 动态数据卷映射表单项 */}
            {renderDynamicVolumeInputs()}

            <div className="lg:flex md:flex hidden justify-start space-x-3">
              <Button
                type="submit"
                variant="surely"
                className="cursor-pointer"
                disabled={loading}
              >
                {t("上传")}
              </Button>

              <SheetClose className="cursor-pointer border border-gray-200/60 bg-gray-200/60 text-black/70 rounded-md text-sm shadow-sm hover:bg-white hover:border-theme-color/85 hover:text-theme-color/85 h-9 px-5 py-2">
                {t("取消")}
              </SheetClose>
            </div>

            {/* 添加小屏幕下的固定按钮组 */}
            <div className="lg:hidden md:hidden flex absolute -top-20 pt-1 right-0  z-50">
                    <Button
                        type="submit"
                        variant="minsure"
                        className="cursor-pointer text-theme-color text-lg font-normal"
                        onClick={handleSubmit(handleUpload)}
                        disabled={loading}
                    >
                        {t('上传')}
                    </Button>
                </div>
          </form>
        </Form>
      </SheetContent>
    </Sheet>
  );
}