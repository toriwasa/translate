#Requires AutoHotkey v2.0

; Ctrl + Cを2回押下することでクリップボードの内容を翻訳する

; グローバル変数の初期化
global tapCount := 0
global lastTapTime := 0

; Ctrl + C のホットキーを設定
; Ctrl + C の入力自体は素通しして実行する
~^c::
{
    ; 現在の時間を取得
    currentTime := A_TickCount

    global tapCount
    global lastTapTime
    ; 起動後初回 OR 前回のタップから一定時間経過していた場合: 押下カウントをリセットして終了する
    if (currentTime - lastTapTime > 2000) {
        tapCount := 1
        ; MsgBox("delayed tap. tapCount: " . tapCount . ", currentTime: " . currentTime . ", lastTapTime: " . lastTapTime)
        ; 現在の時間を保存して次の実行を待機する
        lastTapTime := currentTime
        return
    }
    ; 現在の時間を保存
    lastTapTime := currentTime

    ; 一定時間内の2回目以降の Ctrl + C は押下カウントアップ
    tapCount++
    ; MsgBox("not delayed tap. tapCount: " . tapCount . ", currentTime: " . currentTime . ", lastTapTime: " . lastTapTime)

    ; Ctrl + C が2回以上連続で押下された場合
    if (tapCount >= 2) {
        ; 翻訳ツールの実行ファイルパスを取得する
        translateExePath := A_ScriptDir . "\translate.exe"

        ; このahkファイルのあるディレクトリにtranslate.exeがあることを前提とする
        if (!FileExist(translateExePath)) {
            MsgBox("translate.exeが見つかりませんでした。")
            return
        }

        ; クリップボードの内容を翻訳する処理を実行
        command := translateExePath . " -c"
        Run command
        ; MsgBox("second tap. tapCount: " . tapCount . ", currentTime: " . currentTime . ", lastTapTime: " . lastTapTime)
        ; 翻訳結果のウィンドウが表示されるまで待機して、表示されたらアクティブにする
        if (WinWait("Translate Result", , 3)) {
            WinActivate
        }

        ; 押下カウントをリセット
        tapCount := 0
    }
    return
}
