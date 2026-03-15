import torch

print(torch.cuda.is_available())
print(torch.cuda.get_device_name(0))
print(torch.cuda.mem_get_info())

import os
print(os.path.exists("best_model_stage2.pth"))