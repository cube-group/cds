package qrcode

import (
    "image/png"
    "os"

    "github.com/boombuler/barcode"
    "github.com/boombuler/barcode/qr"
)


// 生成png二维码
// @param text string 文本
// @param filePath string 文件保存路径
// @return width int 图片宽度
// @return height int 图片高度
// @return error
func Png(text string, filePath string, width int, height int) error {
    qrCode, err := qr.Encode(text, qr.M, qr.Auto)
    if err != nil {
        return err
    }

    qrCode, err = barcode.Scale(qrCode, width, height)
    if err != nil {
        return err
    }
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    if err := png.Encode(file, qrCode); err != nil {
        return err
    }
    return nil
}
