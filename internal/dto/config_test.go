package dto

import (
	"reflect"
	"testing"

	agentpb "github.com/ethereal3x/mint-server/api/gen/go/mint_server/agent"
	"github.com/ethereal3x/mint-server/internal/model"
)

// TestConfigToProtoDefaultCapabilities 验证旧文本模型默认补齐文本对话能力
func TestConfigToProtoDefaultCapabilities(t *testing.T) {
	config := &model.ChatModelConfig{ModelType: "gpt-4o-mini"}

	got := ConfigToProto(config)
	want := []string{model.MODEL_CAPABILITY_TEXT_CHAT}
	if !reflect.DeepEqual(got.GetModelCapabilities(), want) {
		t.Fatalf("model capabilities = %v, want %v", got.GetModelCapabilities(), want)
	}
}

// TestConfigToProtoMultimodalCapabilities 验证旧多模态模型默认补齐图片理解能力
func TestConfigToProtoMultimodalCapabilities(t *testing.T) {
	config := &model.ChatModelConfig{
		ModelType:          "gpt-4o",
		SupportsMultimodal: true,
	}

	got := ConfigToProto(config)
	want := []string{model.MODEL_CAPABILITY_TEXT_CHAT, model.MODEL_CAPABILITY_IMAGE_UNDERSTANDING}
	if !reflect.DeepEqual(got.GetModelCapabilities(), want) {
		t.Fatalf("model capabilities = %v, want %v", got.GetModelCapabilities(), want)
	}
}

// TestCreateReqToModelFormatsCapabilities 验证创建请求能力列表会去空和去重后入库
func TestCreateReqToModelFormatsCapabilities(t *testing.T) {
	req := &agentpb.CreateConfigRequest{
		ModelType: "gpt-image-1",
		ModelCapabilities: []string{
			model.MODEL_CAPABILITY_TEXT_TO_IMAGE,
			" ",
			model.MODEL_CAPABILITY_TEXT_TO_IMAGE,
			model.MODEL_CAPABILITY_IMAGE_UNDERSTANDING,
		},
	}

	got := CreateReqToModel(req)
	want := model.MODEL_CAPABILITY_TEXT_TO_IMAGE + "," + model.MODEL_CAPABILITY_IMAGE_UNDERSTANDING
	if got.ModelCapabilities != want {
		t.Fatalf("model capabilities = %q, want %q", got.ModelCapabilities, want)
	}
}
